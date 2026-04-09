package models

type App struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"default:''" json:"description"`
	Icon        string `gorm:"default:''" json:"icon"`
	URL         string `gorm:"default:''" json:"url"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	Order       int    `gorm:"default:0" json:"order"`
}
