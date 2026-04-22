package dto

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
)

const RaiToSqm = 1600.0

type APIResponse struct {
	Success bool            `json:"success"`
	Data    any             `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Error   *ErrorResponse  `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   []ValidationError `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type PlotListQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Search string `form:"search"`
}

type CreatePlotRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	SizeRai   float64 `json:"size_rai" binding:"required,gt=0"`
	TreeCount int     `json:"tree_count" binding:"gte=0"`
	Notes     string  `json:"notes"`
}

type UpdatePlotRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	SizeRai   float64 `json:"size_rai" binding:"required,gt=0"`
	TreeCount int     `json:"tree_count" binding:"gte=0"`
	Notes     string  `json:"notes"`
}

type PlotResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SizeRai   float64   `json:"size_rai"`
	TreeCount int       `json:"tree_count"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeletedPlotResponse struct {
	ID      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}

type PlotStatsQuery struct {
	From string `form:"from"`
	To   string `form:"to"`
}

type PlotStatsResponse struct {
	PlotID      uuid.UUID                    `json:"plot_id"`
	PlotName    string                       `json:"plot_name"`
	Period      PeriodResponse               `json:"period"`
	Sales       PlotStatsSalesResponse       `json:"sales"`
	Maintenance PlotStatsMaintenanceResponse `json:"maintenance"`
}

type PlotStatsSalesResponse struct {
	TotalRevenue      float64                           `json:"total_revenue"`
	TotalWeightKg     float64                           `json:"total_weight_kg"`
	AveragePricePerKg float64                           `json:"average_price_per_kg"`
	TransactionCount  int64                             `json:"transaction_count"`
	GradeBreakdown    map[string]PlotGradeStatsResponse `json:"grade_breakdown"`
}

type PlotGradeStatsResponse struct {
	WeightKg float64 `json:"weight_kg"`
	Revenue  float64 `json:"revenue"`
}

type PlotStatsMaintenanceResponse struct {
	TotalActivities  int64 `json:"total_activities"`
	WateringCount    int64 `json:"watering_count"`
	FertilizingCount int64 `json:"fertilizing_count"`
	PruningCount     int64 `json:"pruning_count"`
	PestControlCount int64 `json:"pest_control_count"`
	HarvestingCount  int64 `json:"harvesting_count"`
}

func NewPaginationMeta(page, limit int, totalItems int64) PaginationMeta {
	totalPages := 0
	if limit > 0 && totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	}

	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1 && totalPages > 0,
	}
}

func NewPlotResponse(plot models.Plot) PlotResponse {
	return PlotResponse{
		ID:        plot.ID,
		Name:      plot.Name,
		SizeRai:   plot.SizeSqm / RaiToSqm,
		TreeCount: plot.TreeCount,
		Notes:     plot.Notes,
		CreatedAt: plot.CreatedAt,
		UpdatedAt: plot.UpdatedAt,
	}
}

func NewPlotResponses(plots []models.Plot) []PlotResponse {
	responses := make([]PlotResponse, len(plots))
	for i, plot := range plots {
		responses[i] = NewPlotResponse(plot)
	}
	return responses
}
