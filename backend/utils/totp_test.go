package utils

import (
	"testing"
	"time"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("Failed to generate TOTP secret: %v", err)
	}

	if len(secret) == 0 {
		t.Error("Secret should not be empty")
	}

	if len(secret) < 32 {
		t.Errorf("Secret should be at least 32 characters, got %d", len(secret))
	}

	for _, c := range secret {
		if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '7')) {
			t.Errorf("Secret contains invalid character: %c", c)
		}
	}
}

func TestGenerateTOTPSecretUniqueness(t *testing.T) {
	secrets := make(map[string]bool)
	for i := 0; i < 100; i++ {
		secret, err := GenerateTOTPSecret()
		if err != nil {
			t.Fatalf("Failed to generate TOTP secret: %v", err)
		}
		if secrets[secret] {
			t.Errorf("Duplicate secret generated: %s", secret)
		}
		secrets[secret] = true
	}
}

func TestGenerateTOTPCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"

	timestamp := int64(1234567890)
	code, err := GenerateTOTPCode(secret, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate TOTP code: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("Code should be 6 digits, got %d", len(code))
	}

	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("Code contains non-digit character: %c", c)
		}
	}
}

func TestGenerateTOTPCodeConsistency(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	timestamp := time.Now().Unix()

	code1, err := GenerateTOTPCode(secret, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate first TOTP code: %v", err)
	}

	code2, err := GenerateTOTPCode(secret, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate second TOTP code: %v", err)
	}

	if code1 != code2 {
		t.Errorf("Same timestamp should produce same code, got %s and %s", code1, code2)
	}
}

func TestGenerateTOTPCodeDifferentSecrets(t *testing.T) {
	secret1 := "JBSWY3DPEHPK3PXP"
	secret2 := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	timestamp := time.Now().Unix()

	code1, err := GenerateTOTPCode(secret1, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate TOTP code for secret1: %v", err)
	}

	code2, err := GenerateTOTPCode(secret2, timestamp)
	if err != nil {
		t.Fatalf("Failed to generate TOTP code for secret2: %v", err)
	}

	if code1 == code2 {
		t.Errorf("Different secrets should produce different codes, got %s for both", code1)
	}
}

func TestVerifyTOTPCode(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("Failed to generate TOTP secret: %v", err)
	}

	code, err := GenerateTOTPCode(secret, time.Now().Unix())
	if err != nil {
		t.Fatalf("Failed to generate TOTP code: %v", err)
	}

	if !VerifyTOTPCode(secret, code) {
		t.Errorf("Valid code should be verified, code: %s", code)
	}
}

func TestVerifyTOTPCodeInvalid(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("Failed to generate TOTP secret: %v", err)
	}

	if VerifyTOTPCode(secret, "000000") {
		t.Error("Invalid code should not be verified")
	}

	if VerifyTOTPCode(secret, "12345") {
		t.Error("5-digit code should not be verified")
	}

	if VerifyTOTPCode(secret, "1234567") {
		t.Error("7-digit code should not be verified")
	}
}

func TestVerifyTOTPCodeTimeWindow(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("Failed to generate TOTP secret: %v", err)
	}

	now := time.Now().Unix()

	code, err := GenerateTOTPCode(secret, now-30)
	if err != nil {
		t.Fatalf("Failed to generate TOTP code: %v", err)
	}
	if !VerifyTOTPCode(secret, code) {
		t.Error("Code from previous time step should be verified within window")
	}

	code, err = GenerateTOTPCode(secret, now+30)
	if err != nil {
		t.Fatalf("Failed to generate TOTP code: %v", err)
	}
	if !VerifyTOTPCode(secret, code) {
		t.Error("Code from next time step should be verified within window")
	}
}

func TestGenerateTOTPURL(t *testing.T) {
	issuer := "TestApp"
	accountName := "user@example.com"
	secret := "JBSWY3DPEHPK3PXP"

	url := GenerateTOTPURL(issuer, accountName, secret)

	expectedPrefix := "otpauth://totp/TestApp:user@example.com?secret=JBSWY3DPEHPK3PXP"
	if len(url) < len(expectedPrefix) || url[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("URL should start with %s, got %s", expectedPrefix, url)
	}

	if !contains(url, "issuer=TestApp") {
		t.Error("URL should contain issuer parameter")
	}

	if !contains(url, "algorithm=SHA1") {
		t.Error("URL should contain algorithm parameter")
	}

	if !contains(url, "digits=6") {
		t.Error("URL should contain digits parameter")
	}

	if !contains(url, "period=30") {
		t.Error("URL should contain period parameter")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
