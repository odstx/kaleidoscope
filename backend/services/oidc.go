package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"kaleidoscope/config"
	"kaleidoscope/models"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type OIDCService struct {
	cfg *config.OIDCConfig
}

func NewOIDCService(cfg *config.OIDCConfig) *OIDCService {
	return &OIDCService{cfg: cfg}
}

func (s *OIDCService) Enabled() bool {
	return s.cfg.Enabled
}

func (s *OIDCService) GetAuthorizationURL(state string) (string, error) {
	if !s.cfg.Enabled {
		return "", errors.New("OIDC is not enabled")
	}

	oauth2Cfg := s.getOAuth2Config()
	authURL := oauth2Cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, nil
}

type OIDCTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type OIDCUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

func (s *OIDCService) ExchangeCode(code string) (*OIDCTokenResponse, error) {
	if !s.cfg.Enabled {
		return nil, errors.New("OIDC is not enabled")
	}

	oauth2Cfg := s.getOAuth2Config()
	token, err := oauth2Cfg.Exchange(nil, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &OIDCTokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int(token.Expiry.Unix()),
		RefreshToken: token.RefreshToken,
		IDToken:      token.Extra("id_token").(string),
	}, nil
}

func (s *OIDCService) GetUserInfo(accessToken string) (*OIDCUserInfo, error) {
	if !s.cfg.Enabled {
		return nil, errors.New("OIDC is not enabled")
	}

	userInfoURL := strings.TrimSuffix(s.cfg.IssuerURL, "/") + "/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo returned status %d", resp.StatusCode)
	}

	var userInfo OIDCUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func (s *OIDCService) FindOrCreateUser(db *gorm.DB, userInfo *OIDCUserInfo) (*models.User, error) {
	var user models.User
	err := db.Where("oidc_subject = ? AND oidc_provider = ?", userInfo.Sub, s.cfg.IssuerURL).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	username := userInfo.Name
	if username == "" {
		username = userInfo.Email
	}
	if username == "" {
		username = userInfo.Sub
	}

	newUser := &models.User{
		UID:          uuid.New().String(),
		Username:     username,
		Email:        userInfo.Email,
		Password:     "",
		OIDCProvider: s.cfg.IssuerURL,
		OIDCSubject:  userInfo.Sub,
	}

	if err := db.Create(newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (s *OIDCService) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.cfg.ClientID,
		ClientSecret: s.cfg.ClientSecret,
		Scopes:       s.cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  strings.TrimSuffix(s.cfg.IssuerURL, "/") + "/authorize",
			TokenURL: strings.TrimSuffix(s.cfg.IssuerURL, "/") + "/token",
		},
		RedirectURL: s.cfg.RedirectURI,
	}
}
