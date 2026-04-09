// Package models defines GORM models for the Durian Farm Management System.
// Production-ready with UUID primary keys, proper relationships, and validation tags.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// ENUMS & CONSTANTS
// =============================================================================

// Grade represents durian quality grades
type Grade string

const (
	GradeA Grade = "A" // Premium quality
	GradeB Grade = "B" // Standard quality
)

// ActivityType represents maintenance activity types
type ActivityType string

const (
	ActivityWatering    ActivityType = "watering"
	ActivityFertilizing ActivityType = "fertilizing"
	ActivityPruning     ActivityType = "pruning"
	ActivityPestControl ActivityType = "pest_control"
	ActivityHarvesting  ActivityType = "harvesting"
)

// =============================================================================
// BASE MODEL
// =============================================================================

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate generates a UUID if not set
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// =============================================================================
// USER MODEL
// =============================================================================

// User represents a farm owner/operator
type User struct {
	BaseModel
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // Never expose in JSON
	FullName     string `gorm:"type:varchar(100);not null" json:"full_name" validate:"required,min=2,max=100"`
	FarmName     string `gorm:"type:varchar(150)" json:"farm_name,omitempty"`
	Phone        string `gorm:"type:varchar(20)" json:"phone,omitempty"`

	// Relationships (has many)
	Plots           []Plot           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"plots,omitempty"`
	Buyers          []Buyer          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"buyers,omitempty"`
	Sales           []Sale           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"sales,omitempty"`
	MaintenanceLogs []MaintenanceLog `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"maintenance_logs,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// =============================================================================
// PLOT MODEL
// =============================================================================

// Plot represents a farm plot/land parcel
type Plot struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	SizeSqm   float64   `gorm:"type:decimal(10,2);not null;check:size_sqm > 0" json:"size_sqm" validate:"required,gt=0"`
	TreeCount int       `gorm:"default:0;check:tree_count >= 0" json:"tree_count" validate:"gte=0"`
	Notes     string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	User            User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Sales           []Sale           `gorm:"foreignKey:PlotID;constraint:OnDelete:RESTRICT" json:"sales,omitempty"`
	MaintenanceLogs []MaintenanceLog `gorm:"foreignKey:PlotID;constraint:OnDelete:RESTRICT" json:"maintenance_logs,omitempty"`
}

// TableName specifies the table name for Plot
func (Plot) TableName() string {
	return "plots"
}

// =============================================================================
// BUYER MODEL
// =============================================================================

// Buyer represents a customer/buyer contact
type Buyer struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name" validate:"required,min=1,max=150"`
	ContactPhone string    `gorm:"type:varchar(20)" json:"contact_phone,omitempty"`
	ContactEmail string    `gorm:"type:varchar(255)" json:"contact_email,omitempty" validate:"omitempty,email"`
	Address      string    `gorm:"type:text" json:"address,omitempty"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	User  User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Sales []Sale `gorm:"foreignKey:BuyerID;constraint:OnDelete:SET NULL" json:"sales,omitempty"`
}

// TableName specifies the table name for Buyer
func (Buyer) TableName() string {
	return "buyers"
}

// =============================================================================
// SALE MODEL
// =============================================================================

// Sale represents a durian sale transaction
type Sale struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	PlotID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"plot_id" validate:"required"`
	BuyerID    *uuid.UUID `gorm:"type:uuid;index" json:"buyer_id,omitempty"` // Nullable
	SaleDate   time.Time  `gorm:"type:date;not null;index" json:"sale_date" validate:"required"`
	Grade      Grade      `gorm:"type:varchar(1);not null;check:grade IN ('A','B')" json:"grade" validate:"required,oneof=A B"`
	WeightKg   float64    `gorm:"type:decimal(10,2);not null;check:weight_kg > 0" json:"weight_kg" validate:"required,gt=0"`
	PricePerKg float64    `gorm:"type:decimal(10,2);not null;check:price_per_kg > 0" json:"price_per_kg" validate:"required,gt=0"`
	TotalPrice float64    `gorm:"type:decimal(12,2);->" json:"total_price"` // Read-only, computed
	Notes      string     `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	User  User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plot  Plot   `gorm:"foreignKey:PlotID" json:"plot,omitempty"`
	Buyer *Buyer `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

// TableName specifies the table name for Sale
func (Sale) TableName() string {
	return "sales"
}

// BeforeCreate computes total_price before inserting
func (s *Sale) BeforeCreate(tx *gorm.DB) error {
	s.TotalPrice = s.WeightKg * s.PricePerKg
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate computes total_price before updating
func (s *Sale) BeforeUpdate(tx *gorm.DB) error {
	s.TotalPrice = s.WeightKg * s.PricePerKg
	return nil
}

// =============================================================================
// MAINTENANCE LOG MODEL
// =============================================================================

// MaintenanceLog represents a farm maintenance activity
type MaintenanceLog struct {
	BaseModel
	UserID          uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	PlotID          uuid.UUID    `gorm:"type:uuid;not null;index" json:"plot_id" validate:"required"`
	ActivityType    ActivityType `gorm:"type:varchar(20);not null;index;check:activity_type IN ('watering','fertilizing','pruning','pest_control','harvesting')" json:"activity_type" validate:"required,oneof=watering fertilizing pruning pest_control harvesting"`
	LogDate         time.Time    `gorm:"type:date;not null;index" json:"log_date" validate:"required"`
	DurationMinutes *int         `gorm:"check:duration_minutes >= 0" json:"duration_minutes,omitempty" validate:"omitempty,gte=0"`
	Quantity        *float64     `gorm:"type:decimal(10,2)" json:"quantity,omitempty" validate:"omitempty,gt=0"`
	QuantityUnit    string       `gorm:"type:varchar(20)" json:"quantity_unit,omitempty"`
	Notes           string       `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plot Plot `gorm:"foreignKey:PlotID" json:"plot,omitempty"`
}

// TableName specifies the table name for MaintenanceLog
func (MaintenanceLog) TableName() string {
	return "maintenance_logs"
}

// =============================================================================
// RESPONSE DTOs (Data Transfer Objects)
// =============================================================================

// UserResponse is the safe response DTO (excludes password)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	FarmName  string    `json:"farm_name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		FarmName:  u.FarmName,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
	}
}

// SaleSummary is used for dashboard aggregations
type SaleSummary struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalWeight  float64 `json:"total_weight_kg"`
	SaleCount    int64   `json:"sale_count"`
	AvgPriceKg   float64 `json:"avg_price_per_kg"`
}

// GradeBreakdown shows grade distribution
type GradeBreakdown struct {
	Grade       Grade   `json:"grade"`
	TotalWeight float64 `json:"total_weight_kg"`
	TotalSales  float64 `json:"total_sales"`
	Percentage  float64 `json:"percentage"`
}

// MonthlyRevenue for chart data
type MonthlyRevenue struct {
	Month    string  `json:"month"` // "2024-07"
	Revenue  float64 `json:"revenue"`
	WeightKg float64 `json:"weight_kg"`
}

// =============================================================================
// MIGRATION FUNCTION
// =============================================================================

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Plot{},
		&Buyer{},
		&Sale{},
		&MaintenanceLog{},
	)
}

// =============================================================================
// SEED DATA
// =============================================================================

// SeedDatabase populates the database with example data
func SeedDatabase(db *gorm.DB) error {
	// Create user
	user := User{
		Email:        "ahmad@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqL9LZU5oIhHD0O1w3Y/1kUVJLvO2", // "password123"
		FullName:     "Ahmad bin Hassan",
		FarmName:     "Ladang Durian Raub",
		Phone:        "+60-12-345-6789",
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Create plots
	plots := []Plot{
		{UserID: user.ID, Name: "Plot A — Hillside", SizeSqm: 2400.00, TreeCount: 48, Notes: "Older trees, best for Grade A"},
		{UserID: user.ID, Name: "Plot B — Valley", SizeSqm: 1800.00, TreeCount: 36, Notes: "Younger trees planted 2021"},
		{UserID: user.ID, Name: "Plot C — Eastern Ridge", SizeSqm: 3200.00, TreeCount: 64, Notes: "Mixed grades"},
		{UserID: user.ID, Name: "Plot D — North Field", SizeSqm: 1200.00, TreeCount: 24, Notes: "Newly planted 2023"},
	}
	if err := db.Create(&plots).Error; err != nil {
		return err
	}

	// Create buyers
	buyers := []Buyer{
		{UserID: user.ID, Name: "Ahmad Traders", ContactPhone: "+60-13-111-2222", ContactEmail: "ahmad.traders@mail.com"},
		{UserID: user.ID, Name: "KL Export Co.", ContactPhone: "+60-3-8888-9999", ContactEmail: "sales@klexport.com"},
		{UserID: user.ID, Name: "Local Market", ContactPhone: "+60-12-555-6666"},
	}
	if err := db.Create(&buyers).Error; err != nil {
		return err
	}

	// Create sales
	sales := []Sale{
		{UserID: user.ID, PlotID: plots[0].ID, BuyerID: &buyers[0].ID, SaleDate: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Grade: GradeA, WeightKg: 320.00, PricePerKg: 28.00},
		{UserID: user.ID, PlotID: plots[2].ID, BuyerID: &buyers[2].ID, SaleDate: time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC), Grade: GradeB, WeightKg: 210.00, PricePerKg: 14.00},
		{UserID: user.ID, PlotID: plots[0].ID, BuyerID: &buyers[1].ID, SaleDate: time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC), Grade: GradeA, WeightKg: 415.00, PricePerKg: 30.00},
		{UserID: user.ID, PlotID: plots[1].ID, BuyerID: &buyers[0].ID, SaleDate: time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC), Grade: GradeA, WeightKg: 180.00, PricePerKg: 27.00},
		{UserID: user.ID, PlotID: plots[2].ID, BuyerID: &buyers[1].ID, SaleDate: time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), Grade: GradeB, WeightKg: 350.00, PricePerKg: 15.00},
		{UserID: user.ID, PlotID: plots[0].ID, BuyerID: &buyers[0].ID, SaleDate: time.Date(2024, 8, 12, 0, 0, 0, 0, time.UTC), Grade: GradeA, WeightKg: 280.00, PricePerKg: 29.00},
	}
	if err := db.Create(&sales).Error; err != nil {
		return err
	}

	// Create maintenance logs
	durationWatering := 45
	durationPestControl := 60
	durationPruning := 90
	quantityFertilizer := 25.00
	quantityPesticide := 2.50

	maintenanceLogs := []MaintenanceLog{
		{UserID: user.ID, PlotID: plots[0].ID, ActivityType: ActivityWatering, LogDate: time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC), DurationMinutes: &durationWatering, Notes: "Soil dry, extra 20 min"},
		{UserID: user.ID, PlotID: plots[1].ID, ActivityType: ActivityFertilizing, LogDate: time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC), Quantity: &quantityFertilizer, QuantityUnit: "kg", Notes: "NPK 15-15-15 blend"},
		{UserID: user.ID, PlotID: plots[2].ID, ActivityType: ActivityPestControl, LogDate: time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC), DurationMinutes: &durationPestControl, Quantity: &quantityPesticide, QuantityUnit: "liters", Notes: "Spotted fruit borers"},
		{UserID: user.ID, PlotID: plots[1].ID, ActivityType: ActivityPruning, LogDate: time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC), DurationMinutes: &durationPruning, Notes: "Removed dead branches"},
		{UserID: user.ID, PlotID: plots[3].ID, ActivityType: ActivityWatering, LogDate: time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), DurationMinutes: &durationWatering, Notes: "New trees need regular watering"},
		{UserID: user.ID, PlotID: plots[0].ID, ActivityType: ActivityHarvesting, LogDate: time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC), Notes: "Harvested for KL Export order"},
	}
	if err := db.Create(&maintenanceLogs).Error; err != nil {
		return err
	}

	return nil
}
