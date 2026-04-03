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

	user, err := uc.userService.Login(req.Email, req.Password)
	if err != nil {
		uc.logger.Error("Login failed", zap.Error(err))
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
	email, _ := c.Get("email")

	uc.logger.Info("Get user info request", zap.Uint("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"email": email,
	})
}
