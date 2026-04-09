package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kaleidoscope/models"
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

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserController handles user-related operations
type UserController struct {
	logger      *zap.Logger
	userService *services.UserService
	oidcService *services.OIDCService
}

// NewUserController creates a new UserController instance
func NewUserController(logger *zap.Logger, userService *services.UserService, oidcService *services.OIDCService) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
		oidcService: oidcService,
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

// GetUserInfo godoc
// @Summary      Get user information
// @Description  Get current user's information including ID, UID, email, and security settings
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "User information retrieved successfully"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Failed to get user info"
// @Router       /users/info [get]
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
		"username":     user.Username,
		"email":        user.Email,
		"totp_enabled": user.TOTPEnabled,
		"hawk_enabled": user.HawkEnabled,
	})
}

// SetupTOTP godoc
// @Summary      Setup TOTP
// @Description  Generate TOTP secret and URL for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  TOTPSetupResponse  "TOTP setup successful"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Failed to setup TOTP"
// @Router       /users/totp/setup [post]
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

// VerifyTOTP godoc
// @Summary      Verify TOTP
// @Description  Verify TOTP code for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Accept       json
// @Produce      json
// @Param        request  body      TOTPVerifyRequest  true  "TOTP verification request"
// @Success      200  {object}  map[string]interface{}  "TOTP verified successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request or TOTP code"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Failed to verify TOTP"
// @Router       /users/totp/verify [post]
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

// EnableTOTP godoc
// @Summary      Enable TOTP
// @Description  Enable TOTP for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "TOTP enabled successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to enable TOTP"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/totp/enable [post]
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

// DisableTOTP godoc
// @Summary      Disable TOTP
// @Description  Disable TOTP for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "TOTP disabled successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to disable TOTP"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/totp/disable [post]
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

// SetupHawk godoc
// @Summary      Setup Hawk
// @Description  Generate Hawk key for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  HawkSetupResponse  "Hawk setup successful"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Failed to setup Hawk"
// @Router       /users/hawk/setup [post]
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

// EnableHawk godoc
// @Summary      Enable Hawk
// @Description  Enable Hawk authentication for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Hawk enabled successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to enable Hawk"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/hawk/enable [post]
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

// DisableHawk godoc
// @Summary      Disable Hawk
// @Description  Disable Hawk authentication for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Hawk disabled successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to disable Hawk"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/hawk/disable [post]
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

type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

// UpdateUsername godoc
// @Summary      Update username
// @Description  Update the username for the current user
// @Tags         users
// @Security     BearerAuth
// @Security     HawkAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateUsernameRequest  true  "Update username request"
// @Success      200  {object}  map[string]interface{}  "Username updated successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/username [put]
func (uc *UserController) UpdateUsername(c *gin.Context) {
	userID := c.GetUint("userID")

	var req UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid update username request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	uc.logger.Info("Update username request", zap.Uint("userID", userID), zap.String("username", req.Username))

	if err := uc.userService.UpdateUsername(userID, req.Username); err != nil {
		uc.logger.Error("Failed to update username", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

// ForgotPassword godoc
// @Summary      Forgot password
// @Description  Send password reset email to user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      ForgotPasswordRequest  true  "Forgot password request"
// @Success      200  {object}  map[string]interface{}  "Password reset email sent"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Router       /users/forgot-password [post]
func (uc *UserController) ForgotPassword(c *gin.Context) {
	uc.logger.Info("Received forgot password request")

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid forgot password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := uc.userService.ForgotPassword(req.Email); err != nil {
		uc.logger.Error("Forgot password failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using token from email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      ResetPasswordRequest  true  "Reset password request"
// @Success      200  {object}  map[string]interface{}  "Password reset successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request or token"
// @Router       /users/reset-password [post]
func (uc *UserController) ResetPassword(c *gin.Context) {
	uc.logger.Info("Received reset password request")

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid reset password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := uc.userService.ResetPassword(req.Token, req.Password); err != nil {
		uc.logger.Error("Reset password failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

type OIDCAuthURLResponse struct {
	URL string `json:"url"`
}

type OIDCCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type OIDCCallbackResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// OIDCLogin godoc
// @Summary      OIDC Login
// @Description  Redirect to OIDC provider for authentication
// @Tags         users
// @Produce      json
// @Success      200  {object}  OIDCAuthURLResponse  "OIDC authorization URL"
// @Router       /users/oidc/login [get]
func (uc *UserController) OIDCLogin(c *gin.Context) {
	if !uc.oidcService.Enabled() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OIDC is not enabled"})
		return
	}

	state := generateState()
	c.Set("oidc_state", state)

	authURL, err := uc.oidcService.GetAuthorizationURL(state)
	if err != nil {
		uc.logger.Error("Failed to get OIDC auth URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authorization URL"})
		return
	}

	c.JSON(http.StatusOK, OIDCAuthURLResponse{URL: authURL})
}

// OIDCCallback godoc
// @Summary      OIDC Callback
// @Description  Handle OIDC provider callback and authenticate user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      OIDCCallbackRequest  true  "OIDC callback request"
// @Success      200  {object}  OIDCCallbackResponse  "Login successful"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      401  {object}  map[string]interface{}  "Authentication failed"
// @Router       /users/oidc/callback [post]
func (uc *UserController) OIDCCallback(c *gin.Context) {
	if !uc.oidcService.Enabled() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OIDC is not enabled"})
		return
	}

	var req OIDCCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Error("Invalid OIDC callback request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tokenResp, err := uc.oidcService.ExchangeCode(req.Code)
	if err != nil {
		uc.logger.Error("Failed to exchange OIDC code", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange code"})
		return
	}

	userInfo, err := uc.oidcService.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		uc.logger.Error("Failed to get OIDC user info", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user info"})
		return
	}

	user, err := uc.oidcService.FindOrCreateUser(uc.userService.GetDB(), userInfo)
	if err != nil {
		uc.logger.Error("Failed to find or create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		uc.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	uc.logger.Info("OIDC login successful", zap.String("email", user.Email))
	c.JSON(http.StatusOK, OIDCCallbackResponse{Token: token, User: user})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
