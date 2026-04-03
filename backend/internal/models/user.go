package models

import "gorm.io/gorm"

// User represents the user model for the database
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
}
