package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

type PlotHandler struct {
	plots *services.PlotService
}

func NewPlotHandler(plots *services.PlotService) *PlotHandler {
	return &PlotHandler{plots: plots}
}

func (h *PlotHandler) RegisterRoutes(rg *gin.RouterGroup) {
	plots := rg.Group("/plots")
	plots.GET("", h.List)
	plots.GET("/:id", h.Get)
	plots.GET("/:id/sales", h.Sales)
	plots.GET("/:id/maintenance", h.Maintenance)
	plots.GET("/:id/stats", h.Stats)
	plots.POST("", h.Create)
	plots.PUT("/:id", h.Update)
	plots.DELETE("/:id", h.Delete)
}

func (h *PlotHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var query dto.PlotListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	plots, meta, err := h.plots.List(c.Request.Context(), userID, query)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    plots,
		Meta:    &meta,
	})
}

func (h *PlotHandler) Get(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	plot, err := h.plots.Get(c.Request.Context(), userID, plotID)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    plot,
	})
}

func (h *PlotHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreatePlotRequest
	if !bindPlotJSON(c, &req) {
		return
	}
	if !validatePlotFields(c, req.Name) {
		return
	}

	plot, err := h.plots.Create(c.Request.Context(), userID, req)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Data:    plot,
	})
}

func (h *PlotHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	var req dto.UpdatePlotRequest
	if !bindPlotJSON(c, &req) {
		return
	}
	if !validatePlotFields(c, req.Name) {
		return
	}

	plot, err := h.plots.Update(c.Request.Context(), userID, plotID, req)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    plot,
	})
}

func (h *PlotHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	deleted, err := h.plots.Delete(c.Request.Context(), userID, plotID)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    deleted,
	})
}

func (h *PlotHandler) Sales(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	var query dto.SaleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	sales, meta, err := h.plots.Sales(c.Request.Context(), userID, plotID, query)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: sales, Meta: &meta})
}

func (h *PlotHandler) Maintenance(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	var query dto.MaintenanceListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	logs, meta, err := h.plots.Maintenance(c.Request.Context(), userID, plotID, query)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: logs, Meta: &meta})
}

func (h *PlotHandler) Stats(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	plotID, ok := parsePlotID(c)
	if !ok {
		return
	}

	var query dto.PlotStatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", nil)
		return
	}

	stats, err := h.plots.Stats(c.Request.Context(), userID, plotID, query)
	if err != nil {
		h.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Data: stats})
}

func (h *PlotHandler) respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrPlotNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Plot not found", nil)
	case errors.Is(err, services.ErrPlotHasDependencies):
		respondError(c, http.StatusConflict, "HAS_DEPENDENCIES", "Cannot delete plot with existing sales or maintenance records", nil)
	case errors.Is(err, services.ErrInvalidSort):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Unsupported sort field", nil)
	case errors.Is(err, services.ErrUserNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "User not found", nil)
	case errors.Is(err, services.ErrInvalidDate):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid date", nil)
	case errors.Is(err, services.ErrInvalidFilter):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid filter", nil)
	case errors.Is(err, services.ErrInvalidGrade):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid grade", nil)
	case errors.Is(err, services.ErrInvalidActivityType):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid activity type", nil)
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
	}
}

func currentUserID(c *gin.Context) (uuid.UUID, bool) {
	rawUserID := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if rawUserID == "" {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing X-User-ID header", nil)
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid X-User-ID header", nil)
		return uuid.Nil, false
	}

	return userID, true
}

func parsePlotID(c *gin.Context) (uuid.UUID, bool) {
	plotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plot id", nil)
		return uuid.Nil, false
	}

	return plotID, true
}

func bindPlotJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", plotValidationDetails(validationErrors))
			return false
		}

		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Malformed JSON request", nil)
		return false
	}

	return true
}

func validatePlotFields(c *gin.Context, name string) bool {
	if strings.TrimSpace(name) != "" {
		return true
	}

	details := []dto.ValidationError{
		{Field: "name", Message: "Name is required"},
	}
	respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details)
	return false
}

func plotValidationDetails(validationErrors validator.ValidationErrors) []dto.ValidationError {
	details := make([]dto.ValidationError, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := plotJSONField(fieldError.Field())
		details = append(details, dto.ValidationError{
			Field:   field,
			Message: plotValidationMessage(field, fieldError.Tag()),
		})
	}
	return details
}

func plotJSONField(field string) string {
	switch field {
	case "Name":
		return "name"
	case "SizeRai":
		return "size_rai"
	case "TreeCount":
		return "tree_count"
	case "Notes":
		return "notes"
	default:
		return strings.ToLower(field)
	}
}

func plotValidationMessage(field, tag string) string {
	switch field {
	case "name":
		if tag == "required" {
			return "Name is required"
		}
		if tag == "max" {
			return "Name must be at most 100 characters"
		}
		return "Name is invalid"
	case "size_rai":
		return "Size must be a positive number"
	case "tree_count":
		return "Tree count must be zero or greater"
	default:
		return "Invalid value"
	}
}

func respondError(c *gin.Context, status int, code, message string, details []dto.ValidationError) {
	c.JSON(status, dto.APIResponse{
		Success: false,
		Error: &dto.ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
