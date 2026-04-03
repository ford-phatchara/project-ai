package database

import (
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the SQLite database connection and auto-migrates schemas
func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto migrate User model
	db.AutoMigrate(&models.User{})

	DB = db
}
