package database

import (
	"fmt"
	"log"

	"shopProject/internal/config"
	"shopProject/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully.")

	// Auto Migration
	err = db.AutoMigrate(
		&models.User{},
		&models.OTPCode{},
		&models.Category{},
		&models.Product{},
		&models.Comment{},
		&models.CartItem{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Database migration completed.")
	DB = db
	return db, nil
}
