package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/services"
)

func parseUUIDParam(c *gin.Context, name, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", message, nil)
		return uuid.Nil, false
	}
	return id, true
}

func bindRequestJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", genericValidationDetails(req, validationErrors))
			return false
		}

		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Malformed JSON request", nil)
		return false
	}

	return true
}

func genericValidationDetails(req any, validationErrors validator.ValidationErrors) []dto.ValidationError {
	details := make([]dto.ValidationError, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := jsonFieldName(req, fieldError.Field())
		details = append(details, dto.ValidationError{
			Field:   field,
			Message: validationMessage(field, fieldError.Tag()),
		})
	}
	return details
}

func jsonFieldName(req any, structField string) string {
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	field, ok := t.FieldByName(structField)
	if !ok {
		return strings.ToLower(structField)
	}

	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return strings.ToLower(structField)
	}

	name := strings.Split(tag, ",")[0]
	if name == "" {
		return strings.ToLower(structField)
	}
	return name
}

func validationMessage(field, tag string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "gt":
		return field + " must be greater than 0"
	case "gte":
		return field + " must be zero or greater"
	case "max":
		return field + " is too long"
	case "oneof":
		return field + " has an unsupported value"
	default:
		return field + " is invalid"
	}
}

func respondDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrBuyerNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Buyer not found", nil)
	case errors.Is(err, services.ErrSaleNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Sale not found", nil)
	case errors.Is(err, services.ErrMaintenanceNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Maintenance log not found", nil)
	case errors.Is(err, services.ErrInvalidRelatedPlot):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Plot not found", nil)
	case errors.Is(err, services.ErrInvalidRelatedBuyer):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "Buyer not found", nil)
	case errors.Is(err, services.ErrInvalidDate):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid date", nil)
	case errors.Is(err, services.ErrInvalidFilter):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid filter", nil)
	case errors.Is(err, services.ErrInvalidGrade):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid grade", nil)
	case errors.Is(err, services.ErrInvalidActivityType):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid activity type", nil)
	case errors.Is(err, services.ErrInvalidSort):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Unsupported sort field", nil)
	case errors.Is(err, services.ErrInvalidGroupBy):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Unsupported group_by value", nil)
	case errors.Is(err, services.ErrInvalidExportFormat):
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "Only csv export is supported", nil)
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
	}
}
