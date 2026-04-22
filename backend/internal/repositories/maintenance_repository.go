package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaintenanceListOptions struct {
	Page          int
	Limit         int
	SortColumn    string
	SortDirection string
	PlotID        *uuid.UUID
	ActivityType  string
	From          *time.Time
	To            *time.Time
}

type MaintenanceActivityCount struct {
	ActivityType string
	Count        int64
}

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) List(ctx context.Context, userID uuid.UUID, opts MaintenanceListOptions) ([]models.MaintenanceLog, int64, error) {
	query := r.applyFilters(r.db.WithContext(ctx).Model(&models.MaintenanceLog{}), userID, opts)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.MaintenanceLog
	err := query.
		Preload("Plot").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.SortColumn},
			Desc:   opts.SortDirection == "desc",
		}).
		Limit(opts.Limit).
		Offset((opts.Page - 1) * opts.Limit).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *MaintenanceRepository) GetByID(ctx context.Context, userID, logID uuid.UUID) (*models.MaintenanceLog, error) {
	var log models.MaintenanceLog
	err := r.db.WithContext(ctx).
		Preload("Plot").
		Where("id = ? AND user_id = ?", logID, userID).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *MaintenanceRepository) Create(ctx context.Context, log *models.MaintenanceLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *MaintenanceRepository) Update(ctx context.Context, log *models.MaintenanceLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *MaintenanceRepository) Delete(ctx context.Context, log *models.MaintenanceLog) error {
	return r.db.WithContext(ctx).Delete(log).Error
}

func (r *MaintenanceRepository) GetPlot(ctx context.Context, userID, plotID uuid.UUID) (*models.Plot, error) {
	var plot models.Plot
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", plotID, userID).
		First(&plot).Error
	if err != nil {
		return nil, err
	}
	return &plot, nil
}

func (r *MaintenanceRepository) Upcoming(ctx context.Context, userID uuid.UUID, from time.Time, limit int) ([]models.MaintenanceLog, error) {
	var logs []models.MaintenanceLog
	err := r.db.WithContext(ctx).
		Preload("Plot").
		Where("user_id = ? AND log_date >= ?", userID, from).
		Order("log_date ASC").
		Order("created_at ASC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *MaintenanceRepository) ActivityCounts(ctx context.Context, userID uuid.UUID, from, to time.Time, plotID *uuid.UUID) ([]MaintenanceActivityCount, error) {
	query := r.db.WithContext(ctx).
		Model(&models.MaintenanceLog{}).
		Select("activity_type, COUNT(*) AS count").
		Where("user_id = ? AND log_date >= ? AND log_date <= ?", userID, from, to)

	if plotID != nil {
		query = query.Where("plot_id = ?", *plotID)
	}

	var counts []MaintenanceActivityCount
	err := query.Group("activity_type").Scan(&counts).Error
	return counts, err
}

func (r *MaintenanceRepository) applyFilters(query *gorm.DB, userID uuid.UUID, opts MaintenanceListOptions) *gorm.DB {
	query = query.Where("user_id = ?", userID)

	if opts.PlotID != nil {
		query = query.Where("plot_id = ?", *opts.PlotID)
	}
	if opts.ActivityType != "" {
		query = query.Where("activity_type = ?", opts.ActivityType)
	}
	if opts.From != nil {
		query = query.Where("log_date >= ?", *opts.From)
	}
	if opts.To != nil {
		query = query.Where("log_date <= ?", *opts.To)
	}

	return query
}
