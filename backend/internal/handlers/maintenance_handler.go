package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

type MaintenanceHandler struct {
	maintenance *services.MaintenanceService
}

func NewMaintenanceHandler(maintenance *services.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{maintenance: maintenance}
}

func (h *MaintenanceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	maintenance := rg.Group("/maintenance")
	maintenance.GET("", h.List)
	maintenance.GET("/upcoming", h.Upcoming)
	maintenance.GET("/:id", h.Get)
	maintenance.POST("", h.Create)
	maintenance.PUT("/:id", h.Update)
	maintenance.DELETE("/:id", h.Delete)
}

// List godoc
// @Summary List maintenance logs
// @Description Returns a paginated list of maintenance logs for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param sort query string false "Sort field, prefix with - for descending"
// @Param plot_id query string false "Plot UUID"
// @Param activity_type query string false "Activity type"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=[]dto.MaintenanceResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance [get]
func (h *MaintenanceHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.MaintenanceListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	logs, meta, err := h.maintenance.List(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: logs, Meta: &meta})
}

// Get godoc
// @Summary Get a maintenance log
// @Description Returns one maintenance log by ID for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Maintenance UUID"
// @Success 200 {object} dto.APIResponse{data=dto.MaintenanceResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance/{id} [get]
func (h *MaintenanceHandler) Get(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	logID, ok := parseUUIDParam(c, "id", "Invalid maintenance id")
	if !ok {
		return
	}

	log, err := h.maintenance.Get(c.Request.Context(), userID, logID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: log})
}

// Create godoc
// @Summary Create a maintenance log
// @Description Creates a maintenance log for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param request body dto.CreateMaintenanceRequest true "Maintenance payload"
// @Success 201 {object} dto.APIResponse{data=dto.MaintenanceResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance [post]
func (h *MaintenanceHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateMaintenanceRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	log, err := h.maintenance.Create(c.Request.Context(), userID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{Success: true, Data: log})
}

// Update godoc
// @Summary Update a maintenance log
// @Description Updates a maintenance log by ID for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Maintenance UUID"
// @Param request body dto.UpdateMaintenanceRequest true "Maintenance payload"
// @Success 200 {object} dto.APIResponse{data=dto.MaintenanceResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 422 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance/{id} [put]
func (h *MaintenanceHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	logID, ok := parseUUIDParam(c, "id", "Invalid maintenance id")
	if !ok {
		return
	}

	var req dto.UpdateMaintenanceRequest
	if !bindRequestJSON(c, &req) {
		return
	}

	log, err := h.maintenance.Update(c.Request.Context(), userID, logID, req)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: log})
}

// Delete godoc
// @Summary Delete a maintenance log
// @Description Deletes a maintenance log by ID for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param id path string true "Maintenance UUID"
// @Success 200 {object} dto.APIResponse{data=dto.DeletedMaintenanceResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance/{id} [delete]
func (h *MaintenanceHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	logID, ok := parseUUIDParam(c, "id", "Invalid maintenance id")
	if !ok {
		return
	}

	deleted, err := h.maintenance.Delete(c.Request.Context(), userID, logID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: deleted})
}

// Upcoming godoc
// @Summary List upcoming maintenance
// @Description Returns upcoming maintenance logs for the current user.
// @Tags maintenance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User UUID"
// @Param limit query int false "Maximum number of logs"
// @Success 200 {object} dto.APIResponse{data=[]dto.MaintenanceResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse}
// @Router /maintenance/upcoming [get]
func (h *MaintenanceHandler) Upcoming(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.MaintenanceUpcomingQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	logs, err := h.maintenance.Upcoming(c.Request.Context(), userID, query)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: logs})
}
