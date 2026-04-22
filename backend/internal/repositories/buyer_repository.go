package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BuyerListOptions struct {
	Page          int
	Limit         int
	Search        string
	SortColumn    string
	SortDirection string
}

type BuyerRepository struct {
	db *gorm.DB
}

func NewBuyerRepository(db *gorm.DB) *BuyerRepository {
	return &BuyerRepository{db: db}
}

func (r *BuyerRepository) List(ctx context.Context, userID uuid.UUID, opts BuyerListOptions) ([]models.Buyer, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Buyer{}).
		Where("user_id = ?", userID)

	if opts.Search != "" {
		query = query.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var buyers []models.Buyer
	err := query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.SortColumn},
			Desc:   opts.SortDirection == "desc",
		}).
		Limit(opts.Limit).
		Offset((opts.Page - 1) * opts.Limit).
		Find(&buyers).Error
	if err != nil {
		return nil, 0, err
	}

	return buyers, total, nil
}

func (r *BuyerRepository) GetByID(ctx context.Context, userID, buyerID uuid.UUID) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", buyerID, userID).
		First(&buyer).Error
	if err != nil {
		return nil, err
	}
	return &buyer, nil
}

func (r *BuyerRepository) Create(ctx context.Context, buyer *models.Buyer) error {
	return r.db.WithContext(ctx).Create(buyer).Error
}

func (r *BuyerRepository) Update(ctx context.Context, buyer *models.Buyer) error {
	return r.db.WithContext(ctx).Save(buyer).Error
}

func (r *BuyerRepository) Delete(ctx context.Context, buyer *models.Buyer) error {
	return r.db.WithContext(ctx).Delete(buyer).Error
}

func (r *BuyerRepository) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Count(&count).Error
	return count > 0, err
}

func (r *BuyerRepository) Sales(ctx context.Context, userID, buyerID uuid.UUID, opts SaleListOptions) ([]models.Sale, int64, error) {
	opts.BuyerID = &buyerID
	return NewSaleRepository(r.db).List(ctx, userID, opts)
}
