package models

import (
	"time"
)

// User represents the user model for the database
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt *time.Time     `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
}
