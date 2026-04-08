package services

import (
	"context"
	"errors"
	"fmt"
	"kaleidoscope/models"
	"kaleidoscope/utils"
	"kaleidoscope/worker"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	client *worker.Client
}

func NewUserService(db *gorm.DB, client *worker.Client) *UserService {
	return &UserService{db: db, client: client}
}

func (s *UserService) GetDB() *gorm.DB {
	return s.db
}

// Register creates a new user with the provided username, email and password
// It validates the input, hashes the password, and persists to database
func (s *UserService) Register(username, email, password string) (*models.User, error) {
	// Input validation
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error while checking existing user: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save to database
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Enqueue welcome email task
	if s.client != nil {
		if err := s.client.EnqueueSendWelcomeEmail(context.Background(), user.ID, user.Username, user.Email); err != nil {
			// Log the error but don't fail the registration - email sending is best effort
			fmt.Printf("Warning: failed to enqueue welcome email: %v\n", err)
		}
	}

	// Remove password from returned user for security
	user.Password = ""
	return user, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error while finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	user.Password = ""
	return &user, nil
}

func (s *UserService) GenerateTOTP(userID uint) (string, string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	secret, err := utils.GenerateTOTPSecret()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	user.TOTPSecret = secret
	user.TOTPVerified = false
	if err := s.db.Save(&user).Error; err != nil {
		return "", "", fmt.Errorf("failed to save TOTP secret: %w", err)
	}

	totpURL := utils.GenerateTOTPURL("Kaleidoscope", user.Email, secret)
	return secret, totpURL, nil
}

func (s *UserService) VerifyTOTP(userID uint, code string) (bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	if user.TOTPSecret == "" {
		return false, errors.New("TOTP not configured for this user")
	}

	if !utils.VerifyTOTPCode(user.TOTPSecret, code) {
		return false, nil
	}

	user.TOTPVerified = true
	if err := s.db.Save(&user).Error; err != nil {
		return false, fmt.Errorf("failed to update TOTP verification status: %w", err)
	}

	return true, nil
}

func (s *UserService) EnableTOTP(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.TOTPSecret == "" {
		return errors.New("TOTP not configured for this user")
	}

	if !user.TOTPVerified {
		return errors.New("TOTP must be verified before enabling")
	}

	user.TOTPEnabled = true
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to enable TOTP: %w", err)
	}

	return nil
}

func (s *UserService) DisableTOTP(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.TOTPSecret = ""
	user.TOTPEnabled = false
	user.TOTPVerified = false
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to disable TOTP: %w", err)
	}

	return nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	user.Password = ""
	return &user, nil
}

func (s *UserService) LoginWithTOTP(email, password, totpCode string) (*models.User, error) {
	user, err := s.Login(email, password)
	if err != nil {
		return nil, err
	}

	if user.TOTPEnabled {
		if totpCode == "" {
			return nil, errors.New("TOTP code required")
		}

		var fullUser models.User
		if err := s.db.Where("email = ?", email).First(&fullUser).Error; err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if !utils.VerifyTOTPCode(fullUser.TOTPSecret, totpCode) {
			return nil, errors.New("invalid TOTP code")
		}
	}

	return user, nil
}

func (s *UserService) GenerateHawkKey(userID uint) (string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	key, err := utils.GenerateHawkKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate Hawk key: %w", err)
	}

	user.HawkKey = key
	user.HawkEnabled = false
	if err := s.db.Save(&user).Error; err != nil {
		return "", fmt.Errorf("failed to save Hawk key: %w", err)
	}

	return key, nil
}

func (s *UserService) EnableHawk(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.HawkKey == "" {
		return errors.New("Hawk key not configured for this user")
	}

	user.HawkEnabled = true
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to enable Hawk: %w", err)
	}

	return nil
}

func (s *UserService) DisableHawk(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.HawkKey = ""
	user.HawkEnabled = false
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to disable Hawk: %w", err)
	}

	return nil
}

func (s *UserService) ForgotPassword(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("database error while finding user: %w", err)
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour).Unix()

	user.ResetToken = token
	user.ResetTokenExpiresAt = expiresAt
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	if s.client != nil {
		if err := s.client.EnqueueSendPasswordResetEmail(context.Background(), user.ID, user.Username, user.Email, token); err != nil {
			fmt.Printf("Warning: failed to enqueue password reset email: %v\n", err)
		}
	}

	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	if token == "" {
		return errors.New("token is required")
	}
	if newPassword == "" {
		return errors.New("password is required")
	}
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var user models.User
	if err := s.db.Where("reset_token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset token")
		}
		return fmt.Errorf("database error while finding user: %w", err)
	}

	if time.Now().Unix() > user.ResetTokenExpiresAt {
		return errors.New("reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiresAt = 0
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
