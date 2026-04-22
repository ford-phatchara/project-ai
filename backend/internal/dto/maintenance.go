package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
)

type MaintenanceListQuery struct {
	Page         int    `form:"page"`
	Limit        int    `form:"limit"`
	Sort         string `form:"sort"`
	PlotID       string `form:"plot_id"`
	ActivityType string `form:"activity_type"`
	From         string `form:"from"`
	To           string `form:"to"`
}

type CreateMaintenanceRequest struct {
	ActivityDate    string    `json:"activity_date" binding:"required"`
	PlotID          uuid.UUID `json:"plot_id" binding:"required"`
	ActivityType    string    `json:"activity_type" binding:"required,oneof=watering fertilizing pruning pest_control harvesting"`
	DurationMinutes *int      `json:"duration_minutes" binding:"omitempty,gte=0"`
	Quantity        *float64  `json:"quantity" binding:"omitempty,gt=0"`
	QuantityUnit    string    `json:"quantity_unit"`
	Notes           string    `json:"notes"`
}

type UpdateMaintenanceRequest struct {
	ActivityDate    string    `json:"activity_date" binding:"required"`
	PlotID          uuid.UUID `json:"plot_id" binding:"required"`
	ActivityType    string    `json:"activity_type" binding:"required,oneof=watering fertilizing pruning pest_control harvesting"`
	DurationMinutes *int      `json:"duration_minutes" binding:"omitempty,gte=0"`
	Quantity        *float64  `json:"quantity" binding:"omitempty,gt=0"`
	QuantityUnit    string    `json:"quantity_unit"`
	Notes           string    `json:"notes"`
}

type MaintenanceUpcomingQuery struct {
	Limit int `form:"limit"`
}

type MaintenanceResponse struct {
	ID              uuid.UUID             `json:"id"`
	ActivityDate    string                `json:"activity_date"`
	Plot            EntitySummaryResponse `json:"plot"`
	ActivityType    models.ActivityType   `json:"activity_type"`
	DurationMinutes *int                  `json:"duration_minutes,omitempty"`
	Quantity        *float64              `json:"quantity,omitempty"`
	QuantityUnit    string                `json:"quantity_unit,omitempty"`
	Notes           string                `json:"notes,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type DeletedMaintenanceResponse struct {
	ID      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}

func NewMaintenanceResponse(log models.MaintenanceLog) MaintenanceResponse {
	return MaintenanceResponse{
		ID:              log.ID,
		ActivityDate:    log.LogDate.Format(DateLayout),
		Plot:            EntitySummaryResponse{ID: log.Plot.ID, Name: log.Plot.Name},
		ActivityType:    log.ActivityType,
		DurationMinutes: log.DurationMinutes,
		Quantity:        log.Quantity,
		QuantityUnit:    log.QuantityUnit,
		Notes:           log.Notes,
		CreatedAt:       log.CreatedAt,
		UpdatedAt:       log.UpdatedAt,
	}
}

func NewMaintenanceResponses(logs []models.MaintenanceLog) []MaintenanceResponse {
	responses := make([]MaintenanceResponse, len(logs))
	for i, log := range logs {
		responses[i] = NewMaintenanceResponse(log)
	}
	return responses
}
