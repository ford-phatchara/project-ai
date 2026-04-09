package database

import (
	"fmt"
	"os"
	"time"

	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the PostgreSQL database connection and auto-migrates schemas
func InitDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	var db *gorm.DB
	var err error

	// Retry connection up to 10 times
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: Failed to connect to database: %v. Retrying in 5s...\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database after retries: %v", err))
	}

	// Enable required extensions
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic(fmt.Sprintf("Failed to create uuid-ossp extension: %v", err))
	}

	// Auto migrate all registered models
	if err := models.AutoMigrate(db); err != nil {
		panic(fmt.Sprintf("Failed to auto migrate database schemas: %v", err))
	}

	DB = db
}
