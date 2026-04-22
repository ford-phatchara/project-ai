package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/dto"
)

var (
	ErrBuyerNotFound       = errors.New("buyer not found")
	ErrSaleNotFound        = errors.New("sale not found")
	ErrMaintenanceNotFound = errors.New("maintenance log not found")
	ErrInvalidDate         = errors.New("invalid date")
	ErrInvalidFilter       = errors.New("invalid filter")
	ErrInvalidGrade        = errors.New("invalid grade")
	ErrInvalidActivityType = errors.New("invalid activity type")
	ErrInvalidGroupBy      = errors.New("invalid group_by")
	ErrInvalidExportFormat = errors.New("invalid export format")
	ErrInvalidRelatedPlot  = errors.New("related plot not found")
	ErrInvalidRelatedBuyer = errors.New("related buyer not found")
)

func normalizePageLimit(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func parseOptionalUUID(value string) (*uuid.UUID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return nil, ErrInvalidFilter
	}
	return &id, nil
}

func parseRequiredDate(value string) (time.Time, error) {
	date, err := time.Parse(dto.DateLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return date, nil
}

func parseOptionalDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	date, err := parseRequiredDate(value)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func defaultDateRange(from, to string, fallbackFrom time.Time) (time.Time, time.Time, error) {
	toDate := time.Now().UTC()
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.UTC)
	fromDate := fallbackFrom

	if strings.TrimSpace(from) != "" {
		parsedFrom, err := parseRequiredDate(from)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		fromDate = parsedFrom
	}

	if strings.TrimSpace(to) != "" {
		parsedTo, err := parseRequiredDate(to)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		toDate = parsedTo
	}

	if fromDate.After(toDate) {
		return time.Time{}, time.Time{}, ErrInvalidDate
	}

	return fromDate, toDate, nil
}

func percentChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return ((current - previous) / previous) * 100
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func validGrade(grade string) bool {
	return grade == "" || grade == "A" || grade == "B"
}

func validActivityType(activityType string) bool {
	switch activityType {
	case "", "watering", "fertilizing", "pruning", "pest_control", "harvesting":
		return true
	default:
		return false
	}
}
