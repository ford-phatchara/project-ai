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
