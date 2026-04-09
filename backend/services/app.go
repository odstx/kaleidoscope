package services

import (
	"kaleidoscope/models"

	"gorm.io/gorm"
)

type AppService struct {
	db *gorm.DB
}

func NewAppService(db *gorm.DB) *AppService {
	return &AppService{db: db}
}

func (as *AppService) GetDB() *gorm.DB {
	return as.db
}

func (as *AppService) GetAllApps() ([]models.App, error) {
	var apps []models.App
	err := as.db.Order("`order` ASC").Find(&apps).Error
	return apps, err
}

func (as *AppService) GetEnabledApps() ([]models.App, error) {
	var apps []models.App
	err := as.db.Where("enabled = ?", true).Order("`order` ASC").Find(&apps).Error
	return apps, err
}

func (as *AppService) GetAppByID(id string) (*models.App, error) {
	var app models.App
	err := as.db.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (as *AppService) CreateApp(app *models.App) (*models.App, error) {
	err := as.db.Create(app).Error
	return app, err
}

func (as *AppService) UpdateApp(id string, app *models.App) (*models.App, error) {
	var existingApp models.App
	if err := as.db.First(&existingApp, id).Error; err != nil {
		return nil, err
	}

	err := as.db.Model(&existingApp).Updates(app).Error
	return &existingApp, err
}

func (as *AppService) DeleteApp(id string) error {
	return as.db.Delete(&models.App{}, id).Error
}
