package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kaleidoscope/config"
	"kaleidoscope/models"
	"kaleidoscope/utils"
)

func HawkAuth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Hawk.Enabled {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Hawk ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization scheme"})
			c.Abort()
			return
		}

		hawkID, err := verifyHawkRequest(c, db, authHeader, cfg.Hawk.TimestampSkewSecs)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("hawkID", hawkID)
		c.Next()
	}
}

func verifyHawkRequest(c *gin.Context, db *gorm.DB, authHeader string, timestampSkewSecs int) (string, error) {
	hawkID, err := utils.ParseAuthorizationHeader(authHeader)
	if err != nil {
		return "", err
	}

	id, ok := hawkID["id"]
	if !ok {
		return "", errors.New("missing hawk id in header")
	}

	var user models.User
	if err := db.Where("uid = ? AND hawk_enabled = ?", id, true).First(&user).Error; err != nil {
		return "", errors.New("user not found or hawk not enabled")
	}

	if user.HawkKey == "" {
		return "", errors.New("hawk key not configured for user")
	}

	host := c.Request.Host
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}

	port := "80"
	if c.Request.TLS != nil {
		port = "443"
	}
	if strings.Contains(c.Request.Host, ":") {
		parts := strings.Split(c.Request.Host, ":")
		port = parts[1]
	}

	uri := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		uri = uri + "?" + c.Request.URL.RawQuery
	}

	verifiedID, err := utils.VerifyHawkAuth(
		authHeader,
		c.Request.Method,
		host,
		port,
		uri,
		user.HawkKey,
		timestampSkewSecs,
	)
	if err != nil {
		return "", err
	}

	c.Set("userID", user.UID)
	c.Set("email", user.Email)

	return verifiedID, nil
}
