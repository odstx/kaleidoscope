package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kaleidoscope/services"
	"kaleidoscope/utils"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TOTPCode string `json:"totp_code"`
}

// TOTPSetupResponse represents the response for TOTP setup
type TOTPSetupResponse struct {
	Secret string `json:"secret"`
	URL    string `json:"url"`
}

// TOTPVerifyRequest represents the request body for TOTP verification
type TOTPVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

// UserController handles user-related operations
type UserController struct {
	logger      *zap.Logger
	userService *services.UserService
}

// NewUserController creates a new UserController instance
func NewUserController(logger *zap.Logger, userService *services.UserService) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
	}
}

// Register registers a new user
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration request"
// @Success      201      {object}  map[string]interface{}  "User registered successfully"
// @Failure      400      {object}  map[string]interface{}  "Invalid request or registration failed"
// @Router       /users/register [post]
func (uc *UserController) Register(c *gin.Context) {
	uc.logger.Info("Received user registration request")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := uc.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		uc.logger.Error("Registration failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info("User registered successfully", zap.String("email", user.Email))
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// Login authenticates a user
// @Summary      Login user
// @Description  Authenticate user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login request"
// @Success      200      {object}  map[string]interface{}  "Login successful"
// @Failure      400      {object}  map[string]interface{}  "Invalid request format"
// @Failure      401      {object}  map[string]interface{}  "Invalid email or password"
// @Router       /users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	uc.logger.Info("Received user login request")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := uc.userService.LoginWithTOTP(req.Email, req.Password, req.TOTPCode)
	if err != nil {
		uc.logger.Error("Login failed", zap.Error(err))
		if err.Error() == "TOTP code required" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "TOTP code required", "totp_required": true})
			return
		}
		if err.Error() == "invalid TOTP code" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		uc.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	uc.logger.Info("User logged in successfully", zap.String("email", user.Email))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": token})
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		uc.logger.Error("Failed to get user info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	uc.logger.Info("Get user info request", zap.Uint("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"uid":          user.UID,
		"email":        user.Email,
		"totp_enabled": user.TOTPEnabled,
		"hawk_enabled": user.HawkEnabled,
	})
}

func (uc *UserController) SetupTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Setup TOTP request", zap.Uint("userID", userID))

	secret, url, err := uc.userService.GenerateTOTP(userID)
	if err != nil {
		uc.logger.Error("Failed to setup TOTP", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup TOTP"})
		return
	}

	c.JSON(http.StatusOK, TOTPSetupResponse{
		Secret: secret,
		URL:    url,
	})
}

func (uc *UserController) VerifyTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	var req TOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid TOTP verify request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	verified, err := uc.userService.VerifyTOTP(userID, req.Code)
	if err != nil {
		uc.logger.Error("Failed to verify TOTP", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify TOTP"})
		return
	}

	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TOTP code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TOTP verified successfully"})
}

func (uc *UserController) EnableTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Enable TOTP request", zap.Uint("userID", userID))

	if err := uc.userService.EnableTOTP(userID); err != nil {
		uc.logger.Error("Failed to enable TOTP", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TOTP enabled successfully"})
}

func (uc *UserController) DisableTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Disable TOTP request", zap.Uint("userID", userID))

	if err := uc.userService.DisableTOTP(userID); err != nil {
		uc.logger.Error("Failed to disable TOTP", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TOTP disabled successfully"})
}

type HawkSetupResponse struct {
	Key string `json:"key"`
}

func (uc *UserController) SetupHawk(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Setup Hawk request", zap.Uint("userID", userID))

	key, err := uc.userService.GenerateHawkKey(userID)
	if err != nil {
		uc.logger.Error("Failed to setup Hawk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup Hawk"})
		return
	}

	c.JSON(http.StatusOK, HawkSetupResponse{Key: key})
}

func (uc *UserController) EnableHawk(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Enable Hawk request", zap.Uint("userID", userID))

	if err := uc.userService.EnableHawk(userID); err != nil {
		uc.logger.Error("Failed to enable Hawk", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hawk enabled successfully"})
}

func (uc *UserController) DisableHawk(c *gin.Context) {
	userID := c.GetUint("userID")

	uc.logger.Info("Disable Hawk request", zap.Uint("userID", userID))

	if err := uc.userService.DisableHawk(userID); err != nil {
		uc.logger.Error("Failed to disable Hawk", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hawk disabled successfully"})
}
