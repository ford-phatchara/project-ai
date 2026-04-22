package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/repositories"
	"gorm.io/gorm"
)

type MaintenanceService struct {
	maintenance *repositories.MaintenanceRepository
}

func NewMaintenanceService(maintenance *repositories.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{maintenance: maintenance}
}

func (s *MaintenanceService) List(ctx context.Context, userID uuid.UUID, query dto.MaintenanceListQuery) ([]dto.MaintenanceResponse, dto.PaginationMeta, error) {
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

func (s *MaintenanceService) Get(ctx context.Context, userID, logID uuid.UUID) (*dto.MaintenanceResponse, error) {
	log, err := s.maintenance.GetByID(ctx, userID, logID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, err
	}

	response := dto.NewMaintenanceResponse(*log)
	return &response, nil
}

func (s *MaintenanceService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateMaintenanceRequest) (*dto.MaintenanceResponse, error) {
	logDate, err := parseRequiredDate(req.ActivityDate)
	if err != nil {
		return nil, err
	}

	plot, err := s.validatePlot(ctx, userID, req.PlotID)
	if err != nil {
		return nil, err
	}

	log := models.MaintenanceLog{
		UserID:          userID,
		PlotID:          req.PlotID,
		ActivityType:    models.ActivityType(req.ActivityType),
		LogDate:         logDate,
		DurationMinutes: req.DurationMinutes,
		Quantity:        req.Quantity,
		QuantityUnit:    strings.TrimSpace(req.QuantityUnit),
		Notes:           strings.TrimSpace(req.Notes),
		Plot:            *plot,
	}

	if err := s.maintenance.Create(ctx, &log); err != nil {
		return nil, err
	}

	response := dto.NewMaintenanceResponse(log)
	return &response, nil
}

func (s *MaintenanceService) Update(ctx context.Context, userID, logID uuid.UUID, req dto.UpdateMaintenanceRequest) (*dto.MaintenanceResponse, error) {
	log, err := s.maintenance.GetByID(ctx, userID, logID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, err
	}

	logDate, err := parseRequiredDate(req.ActivityDate)
	if err != nil {
		return nil, err
	}

	plot, err := s.validatePlot(ctx, userID, req.PlotID)
	if err != nil {
		return nil, err
	}

	log.PlotID = req.PlotID
	log.ActivityType = models.ActivityType(req.ActivityType)
	log.LogDate = logDate
	log.DurationMinutes = req.DurationMinutes
	log.Quantity = req.Quantity
	log.QuantityUnit = strings.TrimSpace(req.QuantityUnit)
	log.Notes = strings.TrimSpace(req.Notes)
	log.Plot = *plot

	if err := s.maintenance.Update(ctx, log); err != nil {
		return nil, err
	}

	response := dto.NewMaintenanceResponse(*log)
	return &response, nil
}

func (s *MaintenanceService) Delete(ctx context.Context, userID, logID uuid.UUID) (*dto.DeletedMaintenanceResponse, error) {
	log, err := s.maintenance.GetByID(ctx, userID, logID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := s.maintenance.Delete(ctx, log); err != nil {
		return nil, err
	}

	return &dto.DeletedMaintenanceResponse{ID: logID, Deleted: true}, nil
}

func (s *MaintenanceService) Upcoming(ctx context.Context, userID uuid.UUID, query dto.MaintenanceUpcomingQuery) ([]dto.MaintenanceResponse, error) {
	limit := query.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	logs, err := s.maintenance.Upcoming(ctx, userID, today, limit)
	if err != nil {
		return nil, err
	}

	return dto.NewMaintenanceResponses(logs), nil
}

func (s *MaintenanceService) validatePlot(ctx context.Context, userID, plotID uuid.UUID) (*models.Plot, error) {
	if plotID == uuid.Nil {
		return nil, ErrInvalidRelatedPlot
	}

	plot, err := s.maintenance.GetPlot(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidRelatedPlot
	}
	if err != nil {
		return nil, err
	}

	return plot, nil
}

func normalizeMaintenanceListOptions(query dto.MaintenanceListQuery) (repositories.MaintenanceListOptions, error) {
	page, limit := normalizePageLimit(query.Page, query.Limit)

	sort := query.Sort
	if sort == "" {
		sort = "-activity_date"
	}

	sortDirection := "asc"
	sortField := sort
	if strings.HasPrefix(sort, "-") {
		sortDirection = "desc"
		sortField = strings.TrimPrefix(sort, "-")
	}

	sortColumns := map[string]string{
		"activity_date": "log_date",
		"log_date":      "log_date",
		"activity_type": "activity_type",
		"created_at":    "created_at",
		"updated_at":    "updated_at",
	}
	sortColumn, ok := sortColumns[sortField]
	if !ok || sortColumn == "" {
		return repositories.MaintenanceListOptions{}, ErrInvalidSort
	}

	if !validActivityType(query.ActivityType) {
		return repositories.MaintenanceListOptions{}, ErrInvalidActivityType
	}

	plotID, err := parseOptionalUUID(query.PlotID)
	if err != nil {
		return repositories.MaintenanceListOptions{}, err
	}
	from, err := parseOptionalDate(query.From)
	if err != nil {
		return repositories.MaintenanceListOptions{}, err
	}
	to, err := parseOptionalDate(query.To)
	if err != nil {
		return repositories.MaintenanceListOptions{}, err
	}
	if from != nil && to != nil && from.After(*to) {
		return repositories.MaintenanceListOptions{}, ErrInvalidDate
	}

	return repositories.MaintenanceListOptions{
		Page:          page,
		Limit:         limit,
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
		PlotID:        plotID,
		ActivityType:  query.ActivityType,
		From:          from,
		To:            to,
	}, nil
}
