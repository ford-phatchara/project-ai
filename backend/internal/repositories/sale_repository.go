package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleListOptions struct {
	Page          int
	Limit         int
	SortColumn    string
	SortDirection string
	PlotID        *uuid.UUID
	Grade         string
	BuyerID       *uuid.UUID
	From          *time.Time
	To            *time.Time
	MinWeight     float64
	MaxWeight     float64
	MinPrice      float64
	MaxPrice      float64
}

type SalesTotals struct {
	Revenue          float64
	WeightKg         float64
	TransactionCount int64
}

type GradeTotals struct {
	Grade    string
	Revenue  float64
	WeightKg float64
}

type TimelineTotals struct {
	Period           string
	Revenue          float64
	WeightKg         float64
	TransactionCount int64
}

type SaleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) List(ctx context.Context, userID uuid.UUID, opts SaleListOptions) ([]models.Sale, int64, error) {
	query := r.applyFilters(r.db.WithContext(ctx).Model(&models.Sale{}), userID, opts)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sales []models.Sale
	err := query.
		Preload("Plot").
		Preload("Buyer").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.SortColumn},
			Desc:   opts.SortDirection == "desc",
		}).
		Limit(opts.Limit).
		Offset((opts.Page - 1) * opts.Limit).
		Find(&sales).Error
	if err != nil {
		return nil, 0, err
	}

	return sales, total, nil
}

func (r *SaleRepository) Export(ctx context.Context, userID uuid.UUID, opts SaleListOptions) ([]models.Sale, error) {
	var sales []models.Sale
	err := r.applyFilters(r.db.WithContext(ctx).Model(&models.Sale{}), userID, opts).
		Preload("Plot").
		Preload("Buyer").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.SortColumn},
			Desc:   opts.SortDirection == "desc",
		}).
		Find(&sales).Error
	return sales, err
}

func (r *SaleRepository) GetByID(ctx context.Context, userID, saleID uuid.UUID) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.WithContext(ctx).
		Preload("Plot").
		Preload("Buyer").
		Where("id = ? AND user_id = ?", saleID, userID).
		First(&sale).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Create(sale).Error
}

func (r *SaleRepository) Update(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Save(sale).Error
}

func (r *SaleRepository) Delete(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Delete(sale).Error
}

func (r *SaleRepository) GetPlot(ctx context.Context, userID, plotID uuid.UUID) (*models.Plot, error) {
	var plot models.Plot
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", plotID, userID).
		First(&plot).Error
	if err != nil {
		return nil, err
	}
	return &plot, nil
}

func (r *SaleRepository) GetBuyer(ctx context.Context, userID, buyerID uuid.UUID) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", buyerID, userID).
		First(&buyer).Error
	if err != nil {
		return nil, err
	}
	return &buyer, nil
}

func (r *SaleRepository) Totals(ctx context.Context, userID uuid.UUID, opts SaleListOptions) (SalesTotals, error) {
	var totals SalesTotals
	err := r.applyFilters(r.db.WithContext(ctx).Model(&models.Sale{}), userID, opts).
		Select(`
			COALESCE(SUM(weight_kg * price_per_kg), 0) AS revenue,
			COALESCE(SUM(weight_kg), 0) AS weight_kg,
			COUNT(*) AS transaction_count
		`).
		Scan(&totals).Error
	return totals, err
}

func (r *SaleRepository) GradeTotals(ctx context.Context, userID uuid.UUID, opts SaleListOptions) ([]GradeTotals, error) {
	var totals []GradeTotals
	err := r.applyFilters(r.db.WithContext(ctx).Model(&models.Sale{}), userID, opts).
		Select(`
			grade,
			COALESCE(SUM(weight_kg * price_per_kg), 0) AS revenue,
			COALESCE(SUM(weight_kg), 0) AS weight_kg
		`).
		Group("grade").
		Order("grade ASC").
		Scan(&totals).Error
	return totals, err
}

func (r *SaleRepository) Timeline(ctx context.Context, userID uuid.UUID, opts SaleListOptions, trunc string, format string) ([]TimelineTotals, error) {
	var timeline []TimelineTotals
	err := r.applyFilters(r.db.WithContext(ctx).Model(&models.Sale{}), userID, opts).
		Select(`
			to_char(date_trunc(?, sale_date), ?) AS period,
			COALESCE(SUM(weight_kg * price_per_kg), 0) AS revenue,
			COALESCE(SUM(weight_kg), 0) AS weight_kg,
			COUNT(*) AS transaction_count
		`, trunc, format).
		Group("period").
		Order("period ASC").
		Scan(&timeline).Error
	return timeline, err
}

func (r *SaleRepository) Recent(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]models.Sale, error) {
	var sales []models.Sale
	err := r.db.WithContext(ctx).
		Preload("Plot").
		Where("user_id = ? AND sale_date >= ? AND sale_date <= ?", userID, from, to).
		Order("sale_date DESC").
		Order("created_at DESC").
		Limit(limit).
		Find(&sales).Error
	return sales, err
}

func (r *SaleRepository) TopPlots(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]struct {
	ID      uuid.UUID
	Name    string
	Revenue float64
}, error) {
	var plots []struct {
		ID      uuid.UUID
		Name    string
		Revenue float64
	}

	err := r.db.WithContext(ctx).
		Table("plots").
		Select("plots.id, plots.name, COALESCE(SUM(sales.weight_kg * sales.price_per_kg), 0) AS revenue").
		Joins("LEFT JOIN sales ON sales.plot_id = plots.id AND sales.deleted_at IS NULL AND sales.sale_date >= ? AND sales.sale_date <= ?", from, to).
		Where("plots.user_id = ? AND plots.deleted_at IS NULL", userID).
		Group("plots.id, plots.name").
		Order("revenue DESC").
		Limit(limit).
		Scan(&plots).Error
	return plots, err
}

func (r *SaleRepository) applyFilters(query *gorm.DB, userID uuid.UUID, opts SaleListOptions) *gorm.DB {
	query = query.Where("user_id = ?", userID)

	if opts.PlotID != nil {
		query = query.Where("plot_id = ?", *opts.PlotID)
	}
	if opts.Grade != "" {
		query = query.Where("grade = ?", opts.Grade)
	}
	if opts.BuyerID != nil {
		query = query.Where("buyer_id = ?", *opts.BuyerID)
	}
	if opts.From != nil {
		query = query.Where("sale_date >= ?", *opts.From)
	}
	if opts.To != nil {
		query = query.Where("sale_date <= ?", *opts.To)
	}
	if opts.MinWeight > 0 {
		query = query.Where("weight_kg >= ?", opts.MinWeight)
	}
	if opts.MaxWeight > 0 {
		query = query.Where("weight_kg <= ?", opts.MaxWeight)
	}
	if opts.MinPrice > 0 {
		query = query.Where("price_per_kg >= ?", opts.MinPrice)
	}
	if opts.MaxPrice > 0 {
		query = query.Where("price_per_kg <= ?", opts.MaxPrice)
	}

	return query
}
