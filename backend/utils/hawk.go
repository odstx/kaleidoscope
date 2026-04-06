package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HawkCredentials struct {
	ID   string
	Key  string
	Algo string
}

type HawkArtifacts struct {
	Method    string
	Host      string
	Port      string
	URI       string
	Hash      string
	Nonce     string
	Ext       string
	App       string
	Dlg       string
	Timestamp int64
	MAC       string
}

func CalculateMAC(algo string, key string, artifacts *HawkArtifacts) (string, error) {
	if algo != "sha256" {
		return "", fmt.Errorf("unsupported algorithm: %s", algo)
	}

	normalized := fmt.Sprintf("hawk.1.header\n%d\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		artifacts.Timestamp,
		artifacts.Nonce,
		artifacts.Method,
		artifacts.URI,
		artifacts.Host,
		artifacts.Port,
		artifacts.Hash,
		artifacts.Ext,
	)

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(normalized))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func ParseAuthorizationHeader(header string) (map[string]string, error) {
	if !strings.HasPrefix(header, "Hawk ") {
		return nil, fmt.Errorf("invalid Hawk header")
	}

	parts := strings.TrimPrefix(header, "Hawk ")
	result := make(map[string]string)

	for _, part := range strings.Split(parts, ", ") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(kv[1], `"`)
		result[key] = value
	}

	return result, nil
}

func VerifyHawkAuth(header string, method string, host string, port string, uri string, key string, timestampSkewSecs int) (string, error) {
	params, err := ParseAuthorizationHeader(header)
	if err != nil {
		return "", err
	}

	id, ok := params["id"]
	if !ok {
		return "", fmt.Errorf("missing id")
	}

	timestampStr, ok := params["ts"]
	if !ok {
		return "", fmt.Errorf("missing timestamp")
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp")
	}

	now := time.Now().Unix()
	skew := now - timestamp
	if skew < 0 {
		skew = -skew
	}
	if int(skew) > timestampSkewSecs {
		return "", fmt.Errorf("timestamp skew too large")
	}

	nonce, ok := params["nonce"]
	if !ok {
		return "", fmt.Errorf("missing nonce")
	}

	mac, ok := params["mac"]
	if !ok {
		return "", fmt.Errorf("missing mac")
	}

	ext := params["ext"]
	hash := params["hash"]

	artifacts := &HawkArtifacts{
		Method:    strings.ToUpper(method),
		Host:      host,
		Port:      port,
		URI:       uri,
		Hash:      hash,
		Nonce:     nonce,
		Ext:       ext,
		Timestamp: timestamp,
	}

	calculatedMAC, err := CalculateMAC("sha256", key, artifacts)
	if err != nil {
		return "", err
	}

	if !hmac.Equal([]byte(mac), []byte(calculatedMAC)) {
		return "", fmt.Errorf("invalid MAC")
	}

	return id, nil
}

func GenerateHawkKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func BuildURI(path string, query string) string {
	if query == "" {
		return path
	}
	return path + "?" + query
}

func GetPortFromURL(u *url.URL) string {
	if u.Port() != "" {
		return u.Port()
	}
	if u.Scheme == "https" {
		return "443"
	}
	return "80"
}
