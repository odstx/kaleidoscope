package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUID   string         `gorm:"index;not null" json:"user_uid"`
	Messages  string         `gorm:"type:json" json:"messages"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Agent) TableName() string {
	return "agents"
}
