package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlotListOptions struct {
	Page          int
	Limit         int
	Search        string
	SortColumn    string
	SortDirection string
}

type PlotRepository struct {
	db *gorm.DB
}

func NewPlotRepository(db *gorm.DB) *PlotRepository {
	return &PlotRepository{db: db}
}

func (r *PlotRepository) List(ctx context.Context, userID uuid.UUID, opts PlotListOptions) ([]models.Plot, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Plot{}).
		Where("user_id = ?", userID)

	if opts.Search != "" {
		query = query.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var plots []models.Plot
	desc := opts.SortDirection == "desc"
	err := query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.SortColumn},
			Desc:   desc,
		}).
		Limit(opts.Limit).
		Offset((opts.Page - 1) * opts.Limit).
		Find(&plots).Error
	if err != nil {
		return nil, 0, err
	}

	return plots, total, nil
}

func (r *PlotRepository) GetByID(ctx context.Context, userID, plotID uuid.UUID) (*models.Plot, error) {
	var plot models.Plot
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", plotID, userID).
		First(&plot).Error
	if err != nil {
		return nil, err
	}
	return &plot, nil
}

func (r *PlotRepository) Create(ctx context.Context, plot *models.Plot) error {
	return r.db.WithContext(ctx).Create(plot).Error
}

func (r *PlotRepository) Update(ctx context.Context, plot *models.Plot) error {
	return r.db.WithContext(ctx).Save(plot).Error
}

func (r *PlotRepository) Delete(ctx context.Context, plot *models.Plot) error {
	return r.db.WithContext(ctx).Delete(plot).Error
}

func (r *PlotRepository) CountDependencies(ctx context.Context, userID, plotID uuid.UUID) (int64, int64, error) {
	var salesCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Sale{}).
		Where("user_id = ? AND plot_id = ?", userID, plotID).
		Count(&salesCount).Error; err != nil {
		return 0, 0, err
	}

	var maintenanceCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.MaintenanceLog{}).
		Where("user_id = ? AND plot_id = ?", userID, plotID).
		Count(&maintenanceCount).Error; err != nil {
		return 0, 0, err
	}

	return salesCount, maintenanceCount, nil
}

func (r *PlotRepository) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Count(&count).Error
	return count > 0, err
}
