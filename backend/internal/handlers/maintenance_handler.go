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
