package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
)

const DateLayout = "2006-01-02"

type SaleListQuery struct {
	Page      int     `form:"page"`
	Limit     int     `form:"limit"`
	Sort      string  `form:"sort"`
	PlotID    string  `form:"plot_id"`
	Grade     string  `form:"grade"`
	BuyerID   string  `form:"buyer_id"`
	From      string  `form:"from"`
	To        string  `form:"to"`
	MinWeight float64 `form:"min_weight"`
	MaxWeight float64 `form:"max_weight"`
	MinPrice  float64 `form:"min_price"`
	MaxPrice  float64 `form:"max_price"`
	HasMinWgt bool
	HasMaxWgt bool
	HasMinPrc bool
	HasMaxPrc bool
}

type CreateSaleRequest struct {
	SaleDate   string     `json:"sale_date" binding:"required"`
	PlotID     uuid.UUID  `json:"plot_id" binding:"required"`
	BuyerID    *uuid.UUID `json:"buyer_id"`
	Grade      string     `json:"grade" binding:"required,oneof=A B"`
	WeightKg   float64    `json:"weight_kg" binding:"required,gt=0"`
	PricePerKg float64    `json:"price_per_kg" binding:"required,gt=0"`
	Notes      string     `json:"notes"`
}

type UpdateSaleRequest struct {
	SaleDate   string     `json:"sale_date" binding:"required"`
	PlotID     uuid.UUID  `json:"plot_id" binding:"required"`
	BuyerID    *uuid.UUID `json:"buyer_id"`
	Grade      string     `json:"grade" binding:"required,oneof=A B"`
	WeightKg   float64    `json:"weight_kg" binding:"required,gt=0"`
	PricePerKg float64    `json:"price_per_kg" binding:"required,gt=0"`
	Notes      string     `json:"notes"`
}

type SaleSummaryQuery struct {
	From    string `form:"from"`
	To      string `form:"to"`
	GroupBy string `form:"group_by"`
	PlotID  string `form:"plot_id"`
}

type SaleExportQuery struct {
	Format string `form:"format"`
	From   string `form:"from"`
	To     string `form:"to"`
	PlotID string `form:"plot_id"`
	Grade  string `form:"grade"`
}

type SaleResponse struct {
	ID         uuid.UUID             `json:"id"`
	SaleDate   string                `json:"sale_date"`
	Plot       EntitySummaryResponse `json:"plot"`
	Buyer      *BuyerSummaryResponse `json:"buyer"`
	Grade      models.Grade          `json:"grade"`
	WeightKg   float64               `json:"weight_kg"`
	PricePerKg float64               `json:"price_per_kg"`
	TotalPrice float64               `json:"total_price"`
	Notes      string                `json:"notes,omitempty"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type EntitySummaryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type DeletedSaleResponse struct {
	ID      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}

type SalesSummaryResponse struct {
	Period   PeriodResponse                  `json:"period"`
	Totals   SalesTotalsResponse             `json:"totals"`
	ByGrade  map[string]GradeSummaryResponse `json:"by_grade"`
	Timeline []SalesTimelinePointResponse    `json:"timeline"`
}

type PeriodResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type SalesTotalsResponse struct {
	Revenue           float64 `json:"revenue"`
	WeightKg          float64 `json:"weight_kg"`
	TransactionCount  int64   `json:"transaction_count"`
	AveragePricePerKg float64 `json:"average_price_per_kg"`
}

type GradeSummaryResponse struct {
	Revenue    float64 `json:"revenue"`
	WeightKg   float64 `json:"weight_kg"`
	Percentage float64 `json:"percentage"`
}

type SalesTimelinePointResponse struct {
	Period           string  `json:"period"`
	Revenue          float64 `json:"revenue"`
	WeightKg         float64 `json:"weight_kg"`
	TransactionCount int64   `json:"transaction_count"`
}

func NewSaleResponse(sale models.Sale) SaleResponse {
	response := SaleResponse{
		ID:         sale.ID,
		SaleDate:   sale.SaleDate.Format(DateLayout),
		Plot:       EntitySummaryResponse{ID: sale.Plot.ID, Name: sale.Plot.Name},
		Grade:      sale.Grade,
		WeightKg:   sale.WeightKg,
		PricePerKg: sale.PricePerKg,
		TotalPrice: sale.WeightKg * sale.PricePerKg,
		Notes:      sale.Notes,
		CreatedAt:  sale.CreatedAt,
		UpdatedAt:  sale.UpdatedAt,
	}

	if sale.Buyer != nil && sale.Buyer.ID != uuid.Nil {
		response.Buyer = NewBuyerSummaryResponse(*sale.Buyer)
	}

	return response
}

func NewSaleResponses(sales []models.Sale) []SaleResponse {
	responses := make([]SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = NewSaleResponse(sale)
	}
	return responses
}
