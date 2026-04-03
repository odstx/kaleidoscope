package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UID      string `gorm:"uniqueIndex;not null" json:"uid"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == "" {
		u.UID = uuid.New().String()
	}
	return nil
}
