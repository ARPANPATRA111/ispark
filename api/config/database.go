package config

import (
	"fmt"
	"log"
	"os"

	"github.com/iips-oss/ispark/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	log.Printf("Connecting to database at %s:%s...", dbHost, dbPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto Migration
	log.Println("Running AutoMigration...")
	err = DB.AutoMigrate(&models.Student{}, &models.OTP{}, &models.Admin{}, &models.Activity{}, &models.Certificate{}, &models.Enrollment{}, &models.SystemSetting{})

	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed.")

	// Safe data migration: backfill coordinator_id from matching admin Name if it is empty/null
	var admins []models.Admin
	if err := DB.Find(&admins).Error; err == nil {
		for _, admin := range admins {
			if err := DB.Model(&models.Activity{}).
				Where("coordinator = ? AND (coordinator_id = ? OR coordinator_id IS NULL)", admin.Name, "").
				Update("coordinator_id", admin.AdminID).Error; err != nil {
				log.Printf("Warning: Failed to backfill coordinator_id for admin %s: %v", admin.AdminID, err)
			}
		}
	}
}
