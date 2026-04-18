package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kaleidoscope/metrics"
	"kaleidoscope/models"
	"kaleidoscope/utils"
	"kaleidoscope/worker"
)

type UserService struct {
	db                      *gorm.DB
	client                  *worker.Client
	resetTokenExpirationHrs int
	logger                  *zap.Logger
}

func NewUserService(db *gorm.DB, client *worker.Client, resetTokenExpirationHrs int, logger *zap.Logger) *UserService {
	return &UserService{db: db, client: client, resetTokenExpirationHrs: resetTokenExpirationHrs, logger: logger}
}

func (s *UserService) GetDB() *gorm.DB {
	return s.db
}

func (s *UserService) IsLockedOut(user *models.User) (bool, int64) {
	if user.LockoutUntil > 0 && time.Now().Unix() < user.LockoutUntil {
		return true, user.LockoutUntil
	}
	return false, 0
}

func (s *UserService) RecordFailedLogin(user *models.User, maxAttempts int, lockoutMins int) error {
	user.FailedLoginAttempts++
	if user.FailedLoginAttempts >= maxAttempts {
		user.LockoutUntil = time.Now().Add(time.Duration(lockoutMins) * time.Minute).Unix()
	}
	return s.db.Save(user).Error
}

func (s *UserService) ResetFailedLogin(user *models.User) error {
	user.FailedLoginAttempts = 0
	user.LockoutUntil = 0
	return s.db.Save(user).Error
}

func (s *UserService) countUsers() int64 {
	var count int64
	s.db.Model(&models.User{}).Count(&count)
	return count
}

func (s *UserService) updateActiveUsersMetric() {
	count := s.countUsers()
	metrics.SetActiveUsers(float64(count))
}

func (s *UserService) Register(username, email, password string) (*models.User, error) {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("register", time.Since(start), err == nil)
	}()

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

	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error while checking existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if s.client != nil {
		if emailErr := s.client.EnqueueSendWelcomeEmail(context.Background(), user.ID, user.Username, user.Email); emailErr != nil {
			fmt.Printf("Warning: failed to enqueue welcome email: %v\n", emailErr)
			metrics.RecordEmailTask("welcome_email", false)
		} else {
			metrics.RecordEmailTask("welcome_email", true)
		}
	}

	s.updateActiveUsersMetric()

	user.Password = ""
	return user, nil
}

func (s *UserService) Login(email, password string, maxAttempts, lockoutMins int) (*models.User, error) {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordAuthOperation("login", err == nil)
		metrics.RecordUserOperation("login", time.Since(start), err == nil)
	}()

	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logSecurityEvent("login_failed", "medium", email, "user_not_found", nil)
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error while finding user: %w", err)
	}

	if locked, until := s.IsLockedOut(&user); locked {
		s.logSecurityEvent("login_failed", "high", email, "account_locked", map[string]interface{}{
			"lockout_until": until,
		})
		return nil, fmt.Errorf("account is locked until %d", until)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.recordFailedLoginAndLog(&user, maxAttempts, lockoutMins)
		s.logSecurityEvent("login_failed", "medium", email, "invalid_password", map[string]interface{}{
			"failed_attempts": user.FailedLoginAttempts + 1,
		})
		return nil, errors.New("invalid email or password")
	}

	_ = s.ResetFailedLogin(&user)

	if s.logger != nil {
		s.logger.Info("User logged in successfully",
			zap.String("email", email),
			zap.Uint("user_id", user.ID),
			zap.String("uid", user.UID),
		)
	}

	user.Password = ""
	return &user, nil
}

func (s *UserService) recordFailedLoginAndLog(user *models.User, maxAttempts, lockoutMins int) {
	_ = s.RecordFailedLogin(user, maxAttempts, lockoutMins)
	
	if s.logger != nil {
		eventType := "failed_login_attempt"
		severity := "medium"
		details := map[string]interface{}{
			"failed_attempts": user.FailedLoginAttempts,
			"max_attempts":    maxAttempts,
		}
		
		if user.FailedLoginAttempts >= maxAttempts {
			eventType = "account_locked"
			severity = "high"
			details["lockout_until"] = user.LockoutUntil
		}
		
		s.logSecurityEvent(eventType, severity, user.Email, "login_failure", details)
	}
}

func (s *UserService) logSecurityEvent(eventType, severity, email, reason string, details map[string]interface{}) {
	if s.logger == nil {
		return
	}
	
	s.logger.Warn("Security event",
		zap.String("event_type", eventType),
		zap.String("severity", severity),
		zap.String("email", email),
		zap.String("reason", reason),
		zap.Any("details", details),
	)
	
	metrics.RecordSecurityEvent(eventType, severity, time.Duration(0))
}

func (s *UserService) GenerateTOTP(userID uint) (string, string, error) {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordTOTPOperation("generate", err == nil)
		metrics.RecordUserOperation("totp_generate", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordTOTPOperation("verify", err == nil)
		metrics.RecordUserOperation("totp_verify", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordTOTPOperation("enable", err == nil)
		metrics.RecordUserOperation("totp_enable", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordTOTPOperation("disable", err == nil)
		metrics.RecordUserOperation("totp_disable", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("get_user", time.Since(start), err == nil)
	}()

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

func (s *UserService) LoginWithTOTP(email, password, totpCode string, maxAttempts, lockoutMins int) (*models.User, error) {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordAuthOperation("login_totp", err == nil)
		metrics.RecordUserOperation("login_totp", time.Since(start), err == nil)
	}()

	user, err := s.Login(email, password, maxAttempts, lockoutMins)
	if err != nil {
		return nil, err
	}

	if user.TOTPEnabled {
		if totpCode == "" {
			s.logSecurityEvent("login_failed", "high", email, "totp_required", nil)
			return nil, errors.New("TOTP code required")
		}

		var fullUser models.User
		if err := s.db.Where("email = ?", email).First(&fullUser).Error; err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if !utils.VerifyTOTPCode(fullUser.TOTPSecret, totpCode) {
			s.logSecurityEvent("login_failed", "high", email, "invalid_totp", nil)
			return nil, errors.New("invalid TOTP code")
		}
	}

	return user, nil
}

func (s *UserService) GenerateHawkKey(userID uint) (string, error) {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("hawk_generate", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("hawk_enable", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("hawk_disable", time.Since(start), err == nil)
	}()

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
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordAuthOperation("forgot_password", err == nil)
		metrics.RecordUserOperation("forgot_password", time.Since(start), err == nil)
	}()

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
	expirationHrs := s.resetTokenExpirationHrs
	if expirationHrs <= 0 {
		expirationHrs = 1
	}
	expiresAt := time.Now().Add(time.Duration(expirationHrs) * time.Hour).Unix()

	user.ResetToken = token
	user.ResetTokenExpiresAt = expiresAt
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	if s.client != nil {
		if emailErr := s.client.EnqueueSendPasswordResetEmail(context.Background(), user.ID, user.Username, user.Email, token); emailErr != nil {
			fmt.Printf("Warning: failed to enqueue password reset email: %v\n", emailErr)
			metrics.RecordEmailTask("password_reset", false)
		} else {
			metrics.RecordEmailTask("password_reset", true)
		}
	}

	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("reset_password", time.Since(start), err == nil)
	}()

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
			s.logSecurityEvent("reset_password_failed", "medium", "", "invalid_token", nil)
			return errors.New("invalid or expired reset token")
		}
		return fmt.Errorf("database error while finding user: %w", err)
	}

	if time.Now().Unix() > user.ResetTokenExpiresAt {
		s.logSecurityEvent("reset_password_failed", "medium", user.Email, "expired_token", nil)
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

	s.logSecurityEvent("password_reset", "medium", user.Email, "success", nil)

	return nil
}

func (s *UserService) UpdateUsername(userID uint, newUsername string) error {
	start := time.Now()
	var err error
	defer func() {
		metrics.RecordUserOperation("update_username", time.Since(start), err == nil)
	}()

	if newUsername == "" {
		return errors.New("username is required")
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Username = newUsername
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}