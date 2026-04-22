package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/repositories"
)

type DashboardService struct {
	sales       *repositories.SaleRepository
	maintenance *repositories.MaintenanceRepository
}

func NewDashboardService(sales *repositories.SaleRepository, maintenance *repositories.MaintenanceRepository) *DashboardService {
	return &DashboardService{sales: sales, maintenance: maintenance}
}

func (s *DashboardService) Overview(ctx context.Context, userID uuid.UUID, query dto.DashboardQuery) (*dto.DashboardResponse, error) {
	from, to, err := defaultDashboardRange(query.From, query.To)
	if err != nil {
		return nil, err
	}

	opts := repositories.SaleListOptions{From: &from, To: &to}
	totals, err := s.sales.Totals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	durationDays := int(to.Sub(from).Hours()/24) + 1
	previousTo := from.AddDate(0, 0, -1)
	previousFrom := previousTo.AddDate(0, 0, -durationDays+1)
	previousTotals, err := s.sales.Totals(ctx, userID, repositories.SaleListOptions{From: &previousFrom, To: &previousTo})
	if err != nil {
		return nil, err
	}

	grades, err := s.sales.GradeTotals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	topPlots, err := s.sales.TopPlots(ctx, userID, from, to, 5)
	if err != nil {
		return nil, err
	}

	recentSales, err := s.sales.Recent(ctx, userID, from, to, 5)
	if err != nil {
		return nil, err
	}

	upcoming, err := s.maintenance.Upcoming(ctx, userID, to, 5)
	if err != nil {
		return nil, err
	}

	averagePrice := 0.0
	if totals.WeightKg > 0 {
		averagePrice = totals.Revenue / totals.WeightKg
	}

	gradeDistribution := map[string]dto.DashboardGradeResponse{
		"A": {},
		"B": {},
	}
	for _, grade := range grades {
		percentage := 0.0
		if totals.WeightKg > 0 {
			percentage = (grade.WeightKg / totals.WeightKg) * 100
		}
		gradeDistribution[grade.Grade] = dto.DashboardGradeResponse{
			Percentage: percentage,
			WeightKg:   grade.WeightKg,
		}
	}

	response := dto.DashboardResponse{
		Period: dto.PeriodResponse{
			From: from.Format(dto.DateLayout),
			To:   to.Format(dto.DateLayout),
		},
		Stats: dto.DashboardStatsResponse{
			TotalRevenue:      totals.Revenue,
			TotalWeightKg:     totals.WeightKg,
			AveragePricePerKg: averagePrice,
			TransactionCount:  totals.TransactionCount,
		},
		Comparison: dto.DashboardComparisonResponse{
			RevenueChangePercent: percentChange(totals.Revenue, previousTotals.Revenue),
			WeightChangePercent:  percentChange(totals.WeightKg, previousTotals.WeightKg),
		},
		GradeDistribution:   gradeDistribution,
		TopPlots:            make([]dto.DashboardTopPlotResponse, len(topPlots)),
		RecentSales:         make([]dto.DashboardRecentSaleResponse, len(recentSales)),
		UpcomingMaintenance: make([]dto.DashboardUpcomingMaintenanceRecord, len(upcoming)),
	}

	for i, plot := range topPlots {
		response.TopPlots[i] = dto.DashboardTopPlotResponse{
			ID:      plot.ID,
			Name:    plot.Name,
			Revenue: plot.Revenue,
		}
	}

	for i, sale := range recentSales {
		response.RecentSales[i] = dto.DashboardRecentSaleResponse{
			ID:         sale.ID,
			SaleDate:   sale.SaleDate.Format(dto.DateLayout),
			PlotName:   sale.Plot.Name,
			Grade:      sale.Grade,
			WeightKg:   sale.WeightKg,
			TotalPrice: sale.WeightKg * sale.PricePerKg,
		}
	}

	for i, log := range upcoming {
		response.UpcomingMaintenance[i] = dto.DashboardUpcomingMaintenanceRecord{
			ID:           log.ID,
			ActivityDate: log.LogDate.Format(dto.DateLayout),
			PlotName:     log.Plot.Name,
			ActivityType: log.ActivityType,
		}
	}

	return &response, nil
}

func (s *DashboardService) Revenue(ctx context.Context, userID uuid.UUID, query dto.DashboardQuery) (*dto.SalesSummaryResponse, error) {
	from, to, err := defaultDashboardRange(query.From, query.To)
	if err != nil {
		return nil, err
	}

	opts := repositories.SaleListOptions{From: &from, To: &to}
	totals, err := s.sales.Totals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}
	grades, err := s.sales.GradeTotals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}
	timeline, err := s.sales.Timeline(ctx, userID, opts, "day", "YYYY-MM-DD")
	if err != nil {
		return nil, err
	}

	response := buildSalesSummaryResponse(from, to, totals, grades, timeline)
	return &response, nil
}

func (s *DashboardService) Activity(ctx context.Context, userID uuid.UUID, query dto.DashboardQuery) (*dto.DashboardActivityResponse, error) {
	from, to, err := defaultDashboardRange(query.From, query.To)
	if err != nil {
		return nil, err
	}

	recentSales, err := s.sales.Recent(ctx, userID, from, to, 10)
	if err != nil {
		return nil, err
	}
	upcoming, err := s.maintenance.Upcoming(ctx, userID, to, 10)
	if err != nil {
		return nil, err
	}

	response := dto.DashboardActivityResponse{
		RecentSales:         make([]dto.DashboardRecentSaleResponse, len(recentSales)),
		UpcomingMaintenance: make([]dto.DashboardUpcomingMaintenanceRecord, len(upcoming)),
	}

	for i, sale := range recentSales {
		response.RecentSales[i] = dto.DashboardRecentSaleResponse{
			ID:         sale.ID,
			SaleDate:   sale.SaleDate.Format(dto.DateLayout),
			PlotName:   sale.Plot.Name,
			Grade:      sale.Grade,
			WeightKg:   sale.WeightKg,
			TotalPrice: sale.WeightKg * sale.PricePerKg,
		}
	}

	for i, log := range upcoming {
		response.UpcomingMaintenance[i] = dto.DashboardUpcomingMaintenanceRecord{
			ID:           log.ID,
			ActivityDate: log.LogDate.Format(dto.DateLayout),
			PlotName:     log.Plot.Name,
			ActivityType: log.ActivityType,
		}
	}

	return &response, nil
}

func defaultDashboardRange(from, to string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return defaultDateRange(from, to, today.AddDate(0, 0, -30))
}
