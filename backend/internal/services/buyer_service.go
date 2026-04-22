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

type BuyerService struct {
	buyers *repositories.BuyerRepository
}

func NewBuyerService(buyers *repositories.BuyerRepository) *BuyerService {
	return &BuyerService{buyers: buyers}
}

func (s *BuyerService) List(ctx context.Context, userID uuid.UUID, query dto.BuyerListQuery) ([]dto.BuyerResponse, dto.PaginationMeta, error) {
	opts, err := normalizeBuyerListOptions(query)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	buyers, total, err := s.buyers.List(ctx, userID, opts)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	meta := dto.NewPaginationMeta(opts.Page, opts.Limit, total)
	return dto.NewBuyerResponses(buyers), meta, nil
}

func (s *BuyerService) Get(ctx context.Context, userID, buyerID uuid.UUID) (*dto.BuyerResponse, error) {
	buyer, err := s.buyers.GetByID(ctx, userID, buyerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBuyerNotFound
	}
	if err != nil {
		return nil, err
	}

	response := dto.NewBuyerResponse(*buyer)
	return &response, nil
}

func (s *BuyerService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateBuyerRequest) (*dto.BuyerResponse, error) {
	exists, err := s.buyers.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	buyer := models.Buyer{
		UserID:       userID,
		Name:         strings.TrimSpace(req.Name),
		ContactPhone: firstNonEmpty(req.ContactPhone, req.Phone),
		ContactEmail: firstNonEmpty(req.ContactEmail, req.Email),
		Address:      strings.TrimSpace(req.Address),
		Notes:        strings.TrimSpace(req.Notes),
	}

	if err := s.buyers.Create(ctx, &buyer); err != nil {
		return nil, err
	}

	response := dto.NewBuyerResponse(buyer)
	return &response, nil
}

func (s *BuyerService) Update(ctx context.Context, userID, buyerID uuid.UUID, req dto.UpdateBuyerRequest) (*dto.BuyerResponse, error) {
	buyer, err := s.buyers.GetByID(ctx, userID, buyerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBuyerNotFound
	}
	if err != nil {
		return nil, err
	}

	buyer.Name = strings.TrimSpace(req.Name)
	buyer.ContactPhone = firstNonEmpty(req.ContactPhone, req.Phone)
	buyer.ContactEmail = firstNonEmpty(req.ContactEmail, req.Email)
	buyer.Address = strings.TrimSpace(req.Address)
	buyer.Notes = strings.TrimSpace(req.Notes)

	if err := s.buyers.Update(ctx, buyer); err != nil {
		return nil, err
	}

	response := dto.NewBuyerResponse(*buyer)
	return &response, nil
}

func (s *BuyerService) Delete(ctx context.Context, userID, buyerID uuid.UUID) (*dto.DeletedBuyerResponse, error) {
	buyer, err := s.buyers.GetByID(ctx, userID, buyerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBuyerNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := s.buyers.Delete(ctx, buyer); err != nil {
		return nil, err
	}

	return &dto.DeletedBuyerResponse{ID: buyerID, Deleted: true}, nil
}

func (s *BuyerService) Sales(ctx context.Context, userID, buyerID uuid.UUID, query dto.SaleListQuery) ([]dto.SaleResponse, dto.PaginationMeta, error) {
	if _, err := s.buyers.GetByID(ctx, userID, buyerID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dto.PaginationMeta{}, ErrBuyerNotFound
	} else if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	opts, err := normalizeSaleListOptions(query)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	sales, total, err := s.buyers.Sales(ctx, userID, buyerID, opts)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	meta := dto.NewPaginationMeta(opts.Page, opts.Limit, total)
	return dto.NewSaleResponses(sales), meta, nil
}

func normalizeBuyerListOptions(query dto.BuyerListQuery) (repositories.BuyerListOptions, error) {
	page, limit := normalizePageLimit(query.Page, query.Limit)

	sort := query.Sort
	if sort == "" {
		sort = "name"
	}

	sortDirection := "asc"
	sortField := sort
	if strings.HasPrefix(sort, "-") {
		sortDirection = "desc"
		sortField = strings.TrimPrefix(sort, "-")
	}

	sortColumns := map[string]string{
		"name":       "name",
		"phone":      "contact_phone",
		"email":      "contact_email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortColumn, ok := sortColumns[sortField]
	if !ok || sortColumn == "" {
		return repositories.BuyerListOptions{}, ErrInvalidSort
	}

	return repositories.BuyerListOptions{
		Page:          page,
		Limit:         limit,
		Search:        strings.TrimSpace(query.Search),
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
	}, nil
}
