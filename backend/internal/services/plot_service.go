package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/repositories"
	"gorm.io/gorm"
)

var (
	ErrPlotNotFound        = errors.New("plot not found")
	ErrPlotHasDependencies = errors.New("plot has dependencies")
	ErrInvalidSort         = errors.New("invalid sort field")
	ErrUserNotFound        = errors.New("user not found")
)

type PlotService struct {
	plots       *repositories.PlotRepository
	sales       *repositories.SaleRepository
	maintenance *repositories.MaintenanceRepository
}

func NewPlotService(plots *repositories.PlotRepository, sales *repositories.SaleRepository, maintenance *repositories.MaintenanceRepository) *PlotService {
	return &PlotService{plots: plots, sales: sales, maintenance: maintenance}
}

func (s *PlotService) List(ctx context.Context, userID uuid.UUID, query dto.PlotListQuery) ([]dto.PlotResponse, dto.PaginationMeta, error) {
	opts, err := normalizePlotListOptions(query)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	plots, total, err := s.plots.List(ctx, userID, opts)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	meta := dto.NewPaginationMeta(opts.Page, opts.Limit, total)
	return dto.NewPlotResponses(plots), meta, nil
}

func (s *PlotService) Get(ctx context.Context, userID, plotID uuid.UUID) (*dto.PlotResponse, error) {
	plot, err := s.plots.GetByID(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlotNotFound
	}
	if err != nil {
		return nil, err
	}

	response := dto.NewPlotResponse(*plot)
	return &response, nil
}

func (s *PlotService) Create(ctx context.Context, userID uuid.UUID, req dto.CreatePlotRequest) (*dto.PlotResponse, error) {
	exists, err := s.plots.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	plot := models.Plot{
		UserID:    userID,
		Name:      strings.TrimSpace(req.Name),
		SizeSqm:   req.SizeRai * dto.RaiToSqm,
		TreeCount: req.TreeCount,
		Notes:     strings.TrimSpace(req.Notes),
	}

	if err := s.plots.Create(ctx, &plot); err != nil {
		return nil, err
	}

	response := dto.NewPlotResponse(plot)
	return &response, nil
}

func (s *PlotService) Update(ctx context.Context, userID, plotID uuid.UUID, req dto.UpdatePlotRequest) (*dto.PlotResponse, error) {
	plot, err := s.plots.GetByID(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlotNotFound
	}
	if err != nil {
		return nil, err
	}

	plot.Name = strings.TrimSpace(req.Name)
	plot.SizeSqm = req.SizeRai * dto.RaiToSqm
	plot.TreeCount = req.TreeCount
	plot.Notes = strings.TrimSpace(req.Notes)

	if err := s.plots.Update(ctx, plot); err != nil {
		return nil, err
	}

	response := dto.NewPlotResponse(*plot)
	return &response, nil
}

func (s *PlotService) Delete(ctx context.Context, userID, plotID uuid.UUID) (*dto.DeletedPlotResponse, error) {
	plot, err := s.plots.GetByID(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlotNotFound
	}
	if err != nil {
		return nil, err
	}

	salesCount, maintenanceCount, err := s.plots.CountDependencies(ctx, userID, plotID)
	if err != nil {
		return nil, err
	}
	if salesCount > 0 || maintenanceCount > 0 {
		return nil, ErrPlotHasDependencies
	}

	if err := s.plots.Delete(ctx, plot); err != nil {
		return nil, err
	}

	return &dto.DeletedPlotResponse{ID: plotID, Deleted: true}, nil
}

func (s *PlotService) Sales(ctx context.Context, userID, plotID uuid.UUID, query dto.SaleListQuery) ([]dto.SaleResponse, dto.PaginationMeta, error) {
	if _, err := s.plots.GetByID(ctx, userID, plotID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dto.PaginationMeta{}, ErrPlotNotFound
	} else if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	query.PlotID = plotID.String()
	opts, err := normalizeSaleListOptions(query)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	sales, total, err := s.sales.List(ctx, userID, opts)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	meta := dto.NewPaginationMeta(opts.Page, opts.Limit, total)
	return dto.NewSaleResponses(sales), meta, nil
}

func (s *PlotService) Maintenance(ctx context.Context, userID, plotID uuid.UUID, query dto.MaintenanceListQuery) ([]dto.MaintenanceResponse, dto.PaginationMeta, error) {
	if _, err := s.plots.GetByID(ctx, userID, plotID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dto.PaginationMeta{}, ErrPlotNotFound
	} else if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	query.PlotID = plotID.String()
	opts, err := normalizeMaintenanceListOptions(query)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	logs, total, err := s.maintenance.List(ctx, userID, opts)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	meta := dto.NewPaginationMeta(opts.Page, opts.Limit, total)
	return dto.NewMaintenanceResponses(logs), meta, nil
}

func (s *PlotService) Stats(ctx context.Context, userID, plotID uuid.UUID, query dto.PlotStatsQuery) (*dto.PlotStatsResponse, error) {
	plot, err := s.plots.GetByID(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlotNotFound
	}
	if err != nil {
		return nil, err
	}

	from, to, err := defaultDashboardRange(query.From, query.To)
	if err != nil {
		return nil, err
	}

	opts := repositories.SaleListOptions{PlotID: &plotID, From: &from, To: &to}
	totals, err := s.sales.Totals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	grades, err := s.sales.GradeTotals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	counts, err := s.maintenance.ActivityCounts(ctx, userID, from, to, &plotID)
	if err != nil {
		return nil, err
	}

	averagePrice := 0.0
	if totals.WeightKg > 0 {
		averagePrice = totals.Revenue / totals.WeightKg
	}

	gradeBreakdown := map[string]dto.PlotGradeStatsResponse{
		"A": {},
		"B": {},
	}
	for _, grade := range grades {
		gradeBreakdown[grade.Grade] = dto.PlotGradeStatsResponse{
			WeightKg: grade.WeightKg,
			Revenue:  grade.Revenue,
		}
	}

	maintenance := dto.PlotStatsMaintenanceResponse{}
	for _, count := range counts {
		maintenance.TotalActivities += count.Count
		switch count.ActivityType {
		case "watering":
			maintenance.WateringCount = count.Count
		case "fertilizing":
			maintenance.FertilizingCount = count.Count
		case "pruning":
			maintenance.PruningCount = count.Count
		case "pest_control":
			maintenance.PestControlCount = count.Count
		case "harvesting":
			maintenance.HarvestingCount = count.Count
		}
	}

	return &dto.PlotStatsResponse{
		PlotID:   plot.ID,
		PlotName: plot.Name,
		Period: dto.PeriodResponse{
			From: from.Format(dto.DateLayout),
			To:   to.Format(dto.DateLayout),
		},
		Sales: dto.PlotStatsSalesResponse{
			TotalRevenue:      totals.Revenue,
			TotalWeightKg:     totals.WeightKg,
			AveragePricePerKg: averagePrice,
			TransactionCount:  totals.TransactionCount,
			GradeBreakdown:    gradeBreakdown,
		},
		Maintenance: maintenance,
	}, nil
}

func normalizePlotListOptions(query dto.PlotListQuery) (repositories.PlotListOptions, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	sort := query.Sort
	if sort == "" {
		sort = "-created_at"
	}

	sortDirection := "asc"
	sortField := sort
	if strings.HasPrefix(sort, "-") {
		sortDirection = "desc"
		sortField = strings.TrimPrefix(sort, "-")
	}

	sortColumns := map[string]string{
		"name":       "name",
		"size_rai":   "size_sqm",
		"tree_count": "tree_count",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortColumn, ok := sortColumns[sortField]
	if !ok || sortColumn == "" {
		return repositories.PlotListOptions{}, ErrInvalidSort
	}

	return repositories.PlotListOptions{
		Page:          page,
		Limit:         limit,
		Search:        strings.TrimSpace(query.Search),
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
	}, nil
}
