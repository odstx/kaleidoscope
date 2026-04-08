package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UID                 string `gorm:"uniqueIndex;not null" json:"uid"`
	Username            string `gorm:"uniqueIndex;not null" json:"username"`
	Email               string `gorm:"uniqueIndex;not null" json:"email"`
	Password            string `gorm:"default:''" json:"password"`
	TOTPSecret          string `gorm:"default:''" json:"-"`
	TOTPEnabled         bool   `gorm:"default:false" json:"totp_enabled"`
	TOTPVerified        bool   `gorm:"default:false" json:"-"`
	HawkKey             string `gorm:"default:''" json:"-"`
	HawkEnabled         bool   `gorm:"default:false" json:"hawk_enabled"`
	ResetToken          string `gorm:"default:''" json:"-"`
	ResetTokenExpiresAt int64  `gorm:"default:0" json:"-"`
	OIDCProvider        string `gorm:"default:''" json:"oidc_provider"`
	OIDCSubject         string `gorm:"default:''" json:"oidc_subject"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == "" {
		u.UID = uuid.New().String()
	}
	return nil
}
