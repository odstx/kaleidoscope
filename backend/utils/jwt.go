package utils

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretBytes     atomic.Value
	jwtExpirationHours atomic.Int64
	jwtOnce            sync.Once
)

func InitJWTSecret(secret string, expirationHours int) {
	jwtOnce.Do(func() {
		jwtSecretBytes.Store([]byte(secret))
		jwtExpirationHours.Store(int64(expirationHours))
	})
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string) (string, error) {
	expirationHours := jwtExpirationHours.Load()
	expiration := time.Duration(expirationHours) * time.Hour
	if expiration == 0 {
		expiration = 24 * time.Hour
	}
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretBytes := jwtSecretBytes.Load().([]byte)
	return token.SignedString(secretBytes)
}

func ParseToken(tokenString string) (*Claims, error) {
	secretBytes := jwtSecretBytes.Load().([]byte)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretBytes, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
