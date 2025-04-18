package database

import (
	"fmt"
	"trendyol-tracker/pkg/config"
	"trendyol-tracker/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schemas
	err = db.AutoMigrate(
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductAttribute{},
		&models.PriceHistory{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
