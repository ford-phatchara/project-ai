package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

type BuyerHandler struct {
	buyers *services.BuyerService
}

func NewBuyerHandler(buyers *services.BuyerService) *BuyerHandler {
	return &BuyerHandler{buyers: buyers}
}

func (h *BuyerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	buyers := rg.Group("/buyers")
	buyers.GET("", h.List)
	buyers.GET("/:id", h.Get)
	buyers.GET("/:id/sales", h.Sales)
	buyers.POST("", h.Create)
	buyers.PUT("/:id", h.Update)
	buyers.DELETE("/:id", h.Delete)
}

// List godoc
// @Summary List buyers
// @Description Returns a paginated list of buyers for the current user.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search term"
// @Param sort query string false "Sort field, prefix with - for descending"
// @Success 200 {object} dto.APIResponse{data=[]dto.BuyerResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers [get]
func (h *BuyerHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.BuyerListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	buyers, meta, err := h.buyers.List(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: buyers, Meta: &meta})
}

// Get godoc
// @Summary Get a buyer
// @Description Returns one buyer by ID for the current user.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Buyer UUID"
// @Success 200 {object} dto.APIResponse{data=dto.BuyerResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers/{id} [get]
func (h *BuyerHandler) Get(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	buyerID, ok := parseUUIDParam(c, "id", "Invalid buyer id")
	if !ok {
		return
	}

	buyer, err := h.buyers.Get(c.Request.Context(), userID, buyerID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: buyer})
}

// Create godoc
// @Summary Create a buyer
// @Description Creates a buyer for the current user.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param request body dto.CreateBuyerRequest true "Buyer payload"
// @Success 201 {object} dto.APIResponse{data=dto.BuyerResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers [post]
func (h *BuyerHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateBuyerRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	buyer, err := h.buyers.Create(c.Request.Context(), userID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{Success: true, Data: buyer})
}

// Update godoc
// @Summary Update a buyer
// @Description Updates a buyer by ID for the current user.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Buyer UUID"
// @Param request body dto.UpdateBuyerRequest true "Buyer payload"
// @Success 200 {object} dto.APIResponse{data=dto.BuyerResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers/{id} [put]
func (h *BuyerHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	buyerID, ok := parseUUIDParam(c, "id", "Invalid buyer id")
	if !ok {
		return
	}

	var req dto.UpdateBuyerRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	buyer, err := h.buyers.Update(c.Request.Context(), userID, buyerID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: buyer})
}

// Delete godoc
// @Summary Delete a buyer
// @Description Deletes a buyer by ID for the current user.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Buyer UUID"
// @Success 200 {object} dto.APIResponse{data=dto.DeletedBuyerResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers/{id} [delete]
func (h *BuyerHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	buyerID, ok := parseUUIDParam(c, "id", "Invalid buyer id")
	if !ok {
		return
	}

	deleted, err := h.buyers.Delete(c.Request.Context(), userID, buyerID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: deleted})
}

// Sales godoc
// @Summary List buyer sales
// @Description Returns paginated sales for a buyer.
// @Tags buyers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Buyer UUID"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param sort query string false "Sort field, prefix with - for descending"
// @Param plot_id query string false "Plot UUID"
// @Param grade query string false "Sale grade"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param min_weight query number false "Minimum weight in kg"
// @Param max_weight query number false "Maximum weight in kg"
// @Param min_price query number false "Minimum price per kg"
// @Param max_price query number false "Maximum price per kg"
// @Success 200 {object} dto.APIResponse{data=[]dto.SaleResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /buyers/{id}/sales [get]
func (h *BuyerHandler) Sales(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	buyerID, ok := parseUUIDParam(c, "id", "Invalid buyer id")
	if !ok {
		return
	}

	var query dto.SaleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	sales, meta, err := h.buyers.Sales(c.Request.Context(), userID, buyerID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: sales, Meta: &meta})
}
