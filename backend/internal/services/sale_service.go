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

type SaleService struct {
	sales *repositories.SaleRepository
}

func NewSaleService(sales *repositories.SaleRepository) *SaleService {
	return &SaleService{sales: sales}
}

func (s *SaleService) List(ctx context.Context, userID uuid.UUID, query dto.SaleListQuery) ([]dto.SaleResponse, dto.PaginationMeta, error) {
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

func (s *SaleService) Get(ctx context.Context, userID, saleID uuid.UUID) (*dto.SaleResponse, error) {
	sale, err := s.sales.GetByID(ctx, userID, saleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSaleNotFound
	}
	if err != nil {
		return nil, err
	}

	response := dto.NewSaleResponse(*sale)
	return &response, nil
}

func (s *SaleService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateSaleRequest) (*dto.SaleResponse, error) {
	saleDate, err := parseRequiredDate(req.SaleDate)
	if err != nil {
		return nil, err
	}

	buyerID := normalizeBuyerID(req.BuyerID)
	plot, buyer, err := s.validateSaleRelations(ctx, userID, req.PlotID, buyerID)
	if err != nil {
		return nil, err
	}

	sale := models.Sale{
		UserID:     userID,
		PlotID:     req.PlotID,
		BuyerID:    buyerID,
		SaleDate:   saleDate,
		Grade:      models.Grade(req.Grade),
		WeightKg:   req.WeightKg,
		PricePerKg: req.PricePerKg,
		TotalPrice: req.WeightKg * req.PricePerKg,
		Notes:      strings.TrimSpace(req.Notes),
		Plot:       *plot,
		Buyer:      buyer,
	}

	if err := s.sales.Create(ctx, &sale); err != nil {
		return nil, err
	}

	response := dto.NewSaleResponse(sale)
	return &response, nil
}

func (s *SaleService) Update(ctx context.Context, userID, saleID uuid.UUID, req dto.UpdateSaleRequest) (*dto.SaleResponse, error) {
	sale, err := s.sales.GetByID(ctx, userID, saleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSaleNotFound
	}
	if err != nil {
		return nil, err
	}

	saleDate, err := parseRequiredDate(req.SaleDate)
	if err != nil {
		return nil, err
	}

	buyerID := normalizeBuyerID(req.BuyerID)
	plot, buyer, err := s.validateSaleRelations(ctx, userID, req.PlotID, buyerID)
	if err != nil {
		return nil, err
	}

	sale.PlotID = req.PlotID
	sale.BuyerID = buyerID
	sale.SaleDate = saleDate
	sale.Grade = models.Grade(req.Grade)
	sale.WeightKg = req.WeightKg
	sale.PricePerKg = req.PricePerKg
	sale.TotalPrice = req.WeightKg * req.PricePerKg
	sale.Notes = strings.TrimSpace(req.Notes)
	sale.Plot = *plot
	sale.Buyer = buyer

	if err := s.sales.Update(ctx, sale); err != nil {
		return nil, err
	}

	response := dto.NewSaleResponse(*sale)
	return &response, nil
}

func (s *SaleService) Delete(ctx context.Context, userID, saleID uuid.UUID) (*dto.DeletedSaleResponse, error) {
	sale, err := s.sales.GetByID(ctx, userID, saleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSaleNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := s.sales.Delete(ctx, sale); err != nil {
		return nil, err
	}

	return &dto.DeletedSaleResponse{ID: saleID, Deleted: true}, nil
}

func (s *SaleService) Summary(ctx context.Context, userID uuid.UUID, query dto.SaleSummaryQuery) (*dto.SalesSummaryResponse, error) {
	from, to, err := defaultSalesSummaryRange(query.From, query.To)
	if err != nil {
		return nil, err
	}

	plotID, err := parseOptionalUUID(query.PlotID)
	if err != nil {
		return nil, err
	}

	groupBy := query.GroupBy
	if groupBy == "" {
		groupBy = "month"
	}
	trunc, format, err := summaryGroup(groupBy)
	if err != nil {
		return nil, err
	}

	opts := repositories.SaleListOptions{From: &from, To: &to, PlotID: plotID}
	totals, err := s.sales.Totals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	grades, err := s.sales.GradeTotals(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	timeline, err := s.sales.Timeline(ctx, userID, opts, trunc, format)
	if err != nil {
		return nil, err
	}

	response := buildSalesSummaryResponse(from, to, totals, grades, timeline)
	return &response, nil
}

func (s *SaleService) Export(ctx context.Context, userID uuid.UUID, query dto.SaleExportQuery) ([]dto.SaleResponse, error) {
	format := strings.TrimSpace(query.Format)
	if format == "" {
		format = "csv"
	}
	if format != "csv" {
		return nil, ErrInvalidExportFormat
	}

	listQuery := dto.SaleListQuery{
		Limit:  100,
		Sort:   "-sale_date",
		PlotID: query.PlotID,
		Grade:  query.Grade,
		From:   query.From,
		To:     query.To,
	}

	opts, err := normalizeSaleListOptions(listQuery)
	if err != nil {
		return nil, err
	}
	opts.Page = 1
	opts.Limit = 0

	sales, err := s.sales.Export(ctx, userID, opts)
	if err != nil {
		return nil, err
	}

	return dto.NewSaleResponses(sales), nil
}

func (s *SaleService) validateSaleRelations(ctx context.Context, userID, plotID uuid.UUID, buyerID *uuid.UUID) (*models.Plot, *models.Buyer, error) {
	if plotID == uuid.Nil {
		return nil, nil, ErrInvalidRelatedPlot
	}

	plot, err := s.sales.GetPlot(ctx, userID, plotID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, ErrInvalidRelatedPlot
	}
	if err != nil {
		return nil, nil, err
	}

	var buyer *models.Buyer
	if buyerID != nil {
		buyer, err = s.sales.GetBuyer(ctx, userID, *buyerID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidRelatedBuyer
		}
		if err != nil {
			return nil, nil, err
		}
	}

	return plot, buyer, nil
}

func normalizeBuyerID(buyerID *uuid.UUID) *uuid.UUID {
	if buyerID == nil || *buyerID == uuid.Nil {
		return nil
	}
	return buyerID
}

func normalizeSaleListOptions(query dto.SaleListQuery) (repositories.SaleListOptions, error) {
	page, limit := normalizePageLimit(query.Page, query.Limit)

	sort := query.Sort
	if sort == "" {
		sort = "-sale_date"
	}

	sortDirection := "asc"
	sortField := sort
	if strings.HasPrefix(sort, "-") {
		sortDirection = "desc"
		sortField = strings.TrimPrefix(sort, "-")
	}

	sortColumns := map[string]string{
		"sale_date":    "sale_date",
		"grade":        "grade",
		"weight_kg":    "weight_kg",
		"price_per_kg": "price_per_kg",
		"total_price":  "total_price",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
	}
	sortColumn, ok := sortColumns[sortField]
	if !ok || sortColumn == "" {
		return repositories.SaleListOptions{}, ErrInvalidSort
	}

	if !validGrade(query.Grade) {
		return repositories.SaleListOptions{}, ErrInvalidGrade
	}

	plotID, err := parseOptionalUUID(query.PlotID)
	if err != nil {
		return repositories.SaleListOptions{}, err
	}
	buyerID, err := parseOptionalUUID(query.BuyerID)
	if err != nil {
		return repositories.SaleListOptions{}, err
	}
	from, err := parseOptionalDate(query.From)
	if err != nil {
		return repositories.SaleListOptions{}, err
	}
	to, err := parseOptionalDate(query.To)
	if err != nil {
		return repositories.SaleListOptions{}, err
	}
	if from != nil && to != nil && from.After(*to) {
		return repositories.SaleListOptions{}, ErrInvalidDate
	}

	if query.MinWeight < 0 || query.MaxWeight < 0 || query.MinPrice < 0 || query.MaxPrice < 0 {
		return repositories.SaleListOptions{}, ErrInvalidFilter
	}
	if query.MaxWeight > 0 && query.MinWeight > query.MaxWeight {
		return repositories.SaleListOptions{}, ErrInvalidFilter
	}
	if query.MaxPrice > 0 && query.MinPrice > query.MaxPrice {
		return repositories.SaleListOptions{}, ErrInvalidFilter
	}

	return repositories.SaleListOptions{
		Page:          page,
		Limit:         limit,
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
		PlotID:        plotID,
		Grade:         query.Grade,
		BuyerID:       buyerID,
		From:          from,
		To:            to,
		MinWeight:     query.MinWeight,
		MaxWeight:     query.MaxWeight,
		MinPrice:      query.MinPrice,
		MaxPrice:      query.MaxPrice,
	}, nil
}

func defaultSalesSummaryRange(from, to string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	return defaultDateRange(from, to, startOfYear)
}

func summaryGroup(groupBy string) (string, string, error) {
	switch groupBy {
	case "day":
		return "day", "YYYY-MM-DD", nil
	case "week":
		return "week", "IYYY-IW", nil
	case "month":
		return "month", "YYYY-MM", nil
	case "year":
		return "year", "YYYY", nil
	default:
		return "", "", ErrInvalidGroupBy
	}
}

func buildSalesSummaryResponse(from, to time.Time, totals repositories.SalesTotals, grades []repositories.GradeTotals, timeline []repositories.TimelineTotals) dto.SalesSummaryResponse {
	averagePrice := 0.0
	if totals.WeightKg > 0 {
		averagePrice = totals.Revenue / totals.WeightKg
	}

	byGrade := map[string]dto.GradeSummaryResponse{
		"A": {},
		"B": {},
	}
	for _, grade := range grades {
		percentage := 0.0
		if totals.WeightKg > 0 {
			percentage = (grade.WeightKg / totals.WeightKg) * 100
		}
		byGrade[grade.Grade] = dto.GradeSummaryResponse{
			Revenue:    grade.Revenue,
			WeightKg:   grade.WeightKg,
			Percentage: percentage,
		}
	}

	points := make([]dto.SalesTimelinePointResponse, len(timeline))
	for i, point := range timeline {
		points[i] = dto.SalesTimelinePointResponse{
			Period:           point.Period,
			Revenue:          point.Revenue,
			WeightKg:         point.WeightKg,
			TransactionCount: point.TransactionCount,
		}
	}

	return dto.SalesSummaryResponse{
		Period: dto.PeriodResponse{
			From: from.Format(dto.DateLayout),
			To:   to.Format(dto.DateLayout),
		},
		Totals: dto.SalesTotalsResponse{
			Revenue:           totals.Revenue,
			WeightKg:          totals.WeightKg,
			TransactionCount:  totals.TransactionCount,
			AveragePricePerKg: averagePrice,
		},
		ByGrade:  byGrade,
		Timeline: points,
	}
}
