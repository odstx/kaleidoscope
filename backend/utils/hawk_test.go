package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestGenerateHawkKey(t *testing.T) {
	key, err := GenerateHawkKey()
	if err != nil {
		t.Fatalf("Failed to generate Hawk key: %v", err)
	}

	if len(key) == 0 {
		t.Error("Key should not be empty")
	}

	decoded, err := base64Decode(key)
	if err != nil {
		t.Errorf("Key should be valid base64: %v", err)
	}

	if len(decoded) != 32 {
		t.Errorf("Decoded key should be 32 bytes, got %d", len(decoded))
	}
}

func base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func TestGenerateHawkKeyUniqueness(t *testing.T) {
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		key, err := GenerateHawkKey()
		if err != nil {
			t.Fatalf("Failed to generate Hawk key: %v", err)
		}
		if keys[key] {
			t.Errorf("Duplicate key generated: %s", key)
		}
		keys[key] = true
	}
}

func TestCalculateMAC(t *testing.T) {
	key := "werxhqm98rp3ngv998sjpsj9s98qjxh"
	artifacts := &HawkArtifacts{
		Method:    "GET",
		Host:      "example.com",
		Port:      "8080",
		URI:       "/resource",
		Hash:      "",
		Nonce:     "U3MOI0fHsBQ",
		Ext:       "",
		Timestamp: 1353822644,
	}

	mac, err := CalculateMAC("sha256", key, artifacts)
	if err != nil {
		t.Fatalf("Failed to calculate MAC: %v", err)
	}

	if len(mac) == 0 {
		t.Error("MAC should not be empty")
	}
}

func TestCalculateMACUnsupportedAlgorithm(t *testing.T) {
	key := "test-key"
	artifacts := &HawkArtifacts{
		Method:    "GET",
		Host:      "example.com",
		Port:      "80",
		URI:       "/",
		Timestamp: time.Now().Unix(),
		Nonce:     "nonce",
	}

	_, err := CalculateMAC("sha1", key, artifacts)
	if err == nil {
		t.Error("Should return error for unsupported algorithm")
	}
}

func TestParseAuthorizationHeader(t *testing.T) {
	header := `Hawk id="dh37fgj492je", ts="1353822644", nonce="U3MOI0fHsBQ", mac="rM6OqC8FhE4bV9jN1sK5dT2gW3xY7zA8bC9dE0fG1hI="`

	params, err := ParseAuthorizationHeader(header)
	if err != nil {
		t.Fatalf("Failed to parse header: %v", err)
	}

	if params["id"] != "dh37fgj492je" {
		t.Errorf("Expected id 'dh37fgj492je', got '%s'", params["id"])
	}

	if params["ts"] != "1353822644" {
		t.Errorf("Expected ts '1353822644', got '%s'", params["ts"])
	}

	if params["nonce"] != "U3MOI0fHsBQ" {
		t.Errorf("Expected nonce 'U3MOI0fHsBQ', got '%s'", params["nonce"])
	}
}

func TestParseAuthorizationHeaderInvalid(t *testing.T) {
	_, err := ParseAuthorizationHeader("Basic abc123")
	if err == nil {
		t.Error("Should return error for non-Hawk header")
	}
}

func TestVerifyHawkAuth(t *testing.T) {
	key := "werxhqm98rp3ngv998sjpsj9s98qjxh"
	timestamp := time.Now().Unix()
	nonce := "U3MOI0fHsBQ"

	artifacts := &HawkArtifacts{
		Method:    "GET",
		Host:      "example.com",
		Port:      "8080",
		URI:       "/resource",
		Timestamp: timestamp,
		Nonce:     nonce,
		Ext:       "",
		Hash:      "",
	}

	mac, err := CalculateMAC("sha256", key, artifacts)
	if err != nil {
		t.Fatalf("Failed to calculate MAC: %v", err)
	}

	header := fmt.Sprintf(`Hawk id="user123", ts="%d", nonce="%s", mac="%s"`, timestamp, nonce, mac)

	id, err := VerifyHawkAuth(header, "GET", "example.com", "8080", "/resource", key, 60)
	if err != nil {
		t.Fatalf("Failed to verify Hawk auth: %v", err)
	}

	if id != "user123" {
		t.Errorf("Expected id 'user123', got '%s'", id)
	}
}

func TestVerifyHawkAuthInvalidMAC(t *testing.T) {
	key := "werxhqm98rp3ngv998sjpsj9s98qjxh"
	timestamp := time.Now().Unix()

	header := fmt.Sprintf(`Hawk id="user123", ts="%d", nonce="nonce", mac="invalidmac"`, timestamp)

	_, err := VerifyHawkAuth(header, "GET", "example.com", "8080", "/resource", key, 60)
	if err == nil {
		t.Error("Should return error for invalid MAC")
	}
}

func TestVerifyHawkAuthTimestampSkew(t *testing.T) {
	key := "werxhqm98rp3ngv998sjpsj9s98qjxh"
	oldTimestamp := time.Now().Unix() - 120

	artifacts := &HawkArtifacts{
		Method:    "GET",
		Host:      "example.com",
		Port:      "8080",
		URI:       "/resource",
		Timestamp: oldTimestamp,
		Nonce:     "nonce",
		Ext:       "",
		Hash:      "",
	}

	mac, _ := CalculateMAC("sha256", key, artifacts)

	header := fmt.Sprintf(`Hawk id="user123", ts="%d", nonce="nonce", mac="%s"`, oldTimestamp, mac)

	_, err := VerifyHawkAuth(header, "GET", "example.com", "8080", "/resource", key, 60)
	if err == nil {
		t.Error("Should return error for timestamp skew too large")
	}
}

func TestVerifyHawkAuthMissingFields(t *testing.T) {
	key := "test-key"

	_, err := VerifyHawkAuth(`Hawk ts="123", nonce="abc", mac="def"`, "GET", "example.com", "80", "/", key, 60)
	if err == nil {
		t.Error("Should return error for missing id")
	}

	_, err = VerifyHawkAuth(`Hawk id="user", nonce="abc", mac="def"`, "GET", "example.com", "80", "/", key, 60)
	if err == nil {
		t.Error("Should return error for missing ts")
	}

	_, err = VerifyHawkAuth(`Hawk id="user", ts="123", mac="def"`, "GET", "example.com", "80", "/", key, 60)
	if err == nil {
		t.Error("Should return error for missing nonce")
	}

	_, err = VerifyHawkAuth(`Hawk id="user", ts="123", nonce="abc"`, "GET", "example.com", "80", "/", key, 60)
	if err == nil {
		t.Error("Should return error for missing mac")
	}
}

func TestBuildURI(t *testing.T) {
	uri := BuildURI("/path", "query=value")
	expected := "/path?query=value"
	if uri != expected {
		t.Errorf("Expected '%s', got '%s'", expected, uri)
	}

	uri = BuildURI("/path", "")
	expected = "/path"
	if uri != expected {
		t.Errorf("Expected '%s', got '%s'", expected, uri)
	}
}

func TestGetPortFromURL(t *testing.T) {
	u, _ := url.Parse("http://example.com:8080/path")
	port := GetPortFromURL(u)
	if port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", port)
	}

	u, _ = url.Parse("http://example.com/path")
	port = GetPortFromURL(u)
	if port != "80" {
		t.Errorf("Expected port '80' for http, got '%s'", port)
	}

	u, _ = url.Parse("https://example.com/path")
	port = GetPortFromURL(u)
	if port != "443" {
		t.Errorf("Expected port '443' for https, got '%s'", port)
	}
}
