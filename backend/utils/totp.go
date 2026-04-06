package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	digits      = 6
	timeStep    = 30
	delayWindow = 1
)

func GenerateTOTPSecret() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	secret := base32.StdEncoding.EncodeToString(bytes)
	return strings.ToUpper(strings.TrimRight(secret, "=")), nil
}

func GenerateTOTPCode(secret string, timestamp int64) (string, error) {
	secret = strings.ToUpper(strings.TrimRight(secret, "="))
	switch len(secret) % 8 {
	case 1:
		secret += "======="
	case 2:
		secret += "======"
	case 3:
		secret += "====="
	case 4:
		secret += "===="
	case 5:
		secret += "==="
	case 6:
		secret += "=="
	case 7:
		secret += "="
	}
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	counter := uint64(timestamp / timeStep)
	code := generateHOTP(key, counter)
	return fmt.Sprintf("%06d", code), nil
}

func VerifyTOTPCode(secret string, code string) bool {
	if len(code) != digits {
		return false
	}

	timestamp := time.Now().Unix()
	for i := -delayWindow; i <= delayWindow; i++ {
		expectedCode, err := GenerateTOTPCode(secret, timestamp+int64(i)*timeStep)
		if err != nil {
			continue
		}
		if expectedCode == code {
			return true
		}
	}
	return false
}

func generateHOTP(key []byte, counter uint64) int {
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, counter)

	h := hmac.New(sha1.New, key)
	h.Write(counterBytes)
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	binaryCode := (int(hash[offset]&0x7f) << 24) |
		(int(hash[offset+1]&0xff) << 16) |
		(int(hash[offset+2]&0xff) << 8) |
		int(hash[offset+3]&0xff)

	return binaryCode % int(math.Pow10(digits))
}

func GenerateTOTPURL(issuer, accountName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=%d&period=%d",
		issuer, accountName, secret, issuer, digits, timeStep)
}
