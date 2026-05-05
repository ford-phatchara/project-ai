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

// Overview godoc
// @Summary Get dashboard overview
// @Description Returns dashboard overview metrics, recent sales, top plots, and upcoming maintenance for the current user.
// @Tags dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.DashboardResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /dashboard [get]
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

// Revenue godoc
// @Summary Get dashboard revenue
// @Description Returns revenue summary analytics for the current user.
// @Tags dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.SalesSummaryResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /dashboard/revenue [get]
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

// Activity godoc
// @Summary Get dashboard activity
// @Description Returns recent sales and upcoming maintenance activity for the current user.
// @Tags dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.DashboardActivityResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /dashboard/activity [get]
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
