package dto

import (
	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
)

type DashboardQuery struct {
	From string `form:"from"`
	To   string `form:"to"`
}

type DashboardResponse struct {
	Period              PeriodResponse                       `json:"period"`
	Stats               DashboardStatsResponse               `json:"stats"`
	Comparison          DashboardComparisonResponse          `json:"comparison"`
	GradeDistribution   map[string]DashboardGradeResponse    `json:"grade_distribution"`
	TopPlots            []DashboardTopPlotResponse           `json:"top_plots"`
	RecentSales         []DashboardRecentSaleResponse        `json:"recent_sales"`
	UpcomingMaintenance []DashboardUpcomingMaintenanceRecord `json:"upcoming_maintenance"`
}

type DashboardStatsResponse struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalWeightKg     float64 `json:"total_weight_kg"`
	AveragePricePerKg float64 `json:"average_price_per_kg"`
	TransactionCount  int64   `json:"transaction_count"`
}

type DashboardComparisonResponse struct {
	RevenueChangePercent float64 `json:"revenue_change_percent"`
	WeightChangePercent  float64 `json:"weight_change_percent"`
}

type DashboardGradeResponse struct {
	Percentage float64 `json:"percentage"`
	WeightKg   float64 `json:"weight_kg"`
}

type DashboardTopPlotResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Revenue float64   `json:"revenue"`
}

type DashboardRecentSaleResponse struct {
	ID         uuid.UUID    `json:"id"`
	SaleDate   string       `json:"sale_date"`
	PlotName   string       `json:"plot_name"`
	Grade      models.Grade `json:"grade"`
	WeightKg   float64      `json:"weight_kg"`
	TotalPrice float64      `json:"total_price"`
}

type DashboardUpcomingMaintenanceRecord struct {
	ID           uuid.UUID           `json:"id"`
	ActivityDate string              `json:"activity_date"`
	PlotName     string              `json:"plot_name"`
	ActivityType models.ActivityType `json:"activity_type"`
}

type DashboardActivityResponse struct {
	RecentSales         []DashboardRecentSaleResponse        `json:"recent_sales"`
	UpcomingMaintenance []DashboardUpcomingMaintenanceRecord `json:"upcoming_maintenance"`
}
