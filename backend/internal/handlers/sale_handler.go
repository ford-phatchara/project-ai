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

// List godoc
// @Summary List sales
// @Description Returns a paginated list of sales for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param sort query string false "Sort field, prefix with - for descending"
// @Param plot_id query string false "Plot UUID"
// @Param grade query string false "Sale grade"
// @Param buyer_id query string false "Buyer UUID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param min_weight query number false "Minimum weight in kg"
// @Param max_weight query number false "Maximum weight in kg"
// @Param min_price query number false "Minimum price per kg"
// @Param max_price query number false "Maximum price per kg"
// @Success 200 {object} dto.APIResponse{data=[]dto.SaleResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales [get]
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

// Get godoc
// @Summary Get a sale
// @Description Returns one sale by ID for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Sale UUID"
// @Success 200 {object} dto.APIResponse{data=dto.SaleResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales/{id} [get]
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

// Create godoc
// @Summary Create a sale
// @Description Creates a sale for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param request body dto.CreateSaleRequest true "Sale payload"
// @Success 201 {object} dto.APIResponse{data=dto.SaleResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales [post]
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

// Update godoc
// @Summary Update a sale
// @Description Updates a sale by ID for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Sale UUID"
// @Param request body dto.UpdateSaleRequest true "Sale payload"
// @Success 200 {object} dto.APIResponse{data=dto.SaleResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales/{id} [put]
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

// Delete godoc
// @Summary Delete a sale
// @Description Deletes a sale by ID for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Sale UUID"
// @Success 200 {object} dto.APIResponse{data=dto.DeletedSaleResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales/{id} [delete]
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

// Summary godoc
// @Summary Get sales summary
// @Description Returns sales totals, grade breakdown, and timeline data for the current user.
// @Tags sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param group_by query string false "Grouping period"
// @Param plot_id query string false "Plot UUID"
// @Success 200 {object} dto.APIResponse{data=dto.SalesSummaryResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales/summary [get]
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

// Export godoc
// @Summary Export sales
// @Description Exports sales for the current user as a CSV file.
// @Tags sales
// @Accept json
// @Produce text/csv
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param format query string false "Export format, only csv is supported"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param plot_id query string false "Plot UUID"
// @Param grade query string false "Sale grade"
// @Success 200 {file} file "CSV file"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /sales/export [get]
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
