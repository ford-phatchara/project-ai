package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

type DashboardHandler struct {
	dashboard *services.DashboardService
}

func NewDashboardHandler(dashboard *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard}
}

func (h *DashboardHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dashboard := rg.Group("/dashboard")
	dashboard.GET("", h.Overview)
	dashboard.GET("/revenue", h.Revenue)
	dashboard.GET("/activity", h.Activity)
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	overview, err := h.dashboard.Overview(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: overview})
}

func (h *DashboardHandler) Revenue(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	revenue, err := h.dashboard.Revenue(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: revenue})
}

func (h *DashboardHandler) Activity(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	activity, err := h.dashboard.Activity(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: activity})
}
