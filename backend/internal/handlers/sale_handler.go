package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

type SaleHandler struct {
	sales *services.SaleService
}

func NewSaleHandler(sales *services.SaleService) *SaleHandler {
	return &SaleHandler{sales: sales}
}

func (h *SaleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	sales := rg.Group("/sales")
	sales.GET("", h.List)
	sales.GET("/summary", h.Summary)
	sales.GET("/export", h.Export)
	sales.GET("/:id", h.Get)
	sales.POST("", h.Create)
	sales.PUT("/:id", h.Update)
	sales.DELETE("/:id", h.Delete)
}

func (h *SaleHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.SaleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	sales, meta, err := h.sales.List(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: sales, Meta: &meta})
}

func (h *SaleHandler) Get(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	saleID, ok := parseUUIDParam(c, "id", "Invalid sale id")
	if !ok {
		return
	}

	sale, err := h.sales.Get(c.Request.Context(), userID, saleID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: sale})
}

func (h *SaleHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateSaleRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	sale, err := h.sales.Create(c.Request.Context(), userID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{Success: true, Data: sale})
}

func (h *SaleHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	saleID, ok := parseUUIDParam(c, "id", "Invalid sale id")
	if !ok {
		return
	}

	var req dto.UpdateSaleRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	sale, err := h.sales.Update(c.Request.Context(), userID, saleID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: sale})
}

func (h *SaleHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	saleID, ok := parseUUIDParam(c, "id", "Invalid sale id")
	if !ok {
		return
	}

	deleted, err := h.sales.Delete(c.Request.Context(), userID, saleID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: deleted})
}

func (h *SaleHandler) Summary(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.SaleSummaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	summary, err := h.sales.Summary(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: summary})
}

func (h *SaleHandler) Export(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.SaleExportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	sales, err := h.sales.Export(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="sales_export.csv"`)

	writer := csv.NewWriter(c.Writer)
	_ = writer.Write([]string{"sale_date", "plot_name", "buyer_name", "grade", "weight_kg", "price_per_kg", "total_price", "notes"})
	for _, sale := range sales {
		buyerName := ""
		if sale.Buyer != nil {
			buyerName = sale.Buyer.Name
		}
		_ = writer.Write([]string{
			sale.SaleDate,
			sale.Plot.Name,
			buyerName,
			string(sale.Grade),
			fmt.Sprintf("%.2f", sale.WeightKg),
			fmt.Sprintf("%.2f", sale.PricePerKg),
			fmt.Sprintf("%.2f", sale.TotalPrice),
			sale.Notes,
		})
	}
	writer.Flush()
}
