package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
)

type BuyerListQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
}

type CreateBuyerRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=150"`
	ContactName  string `json:"contact_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"omitempty,email"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email"`
	Address      string `json:"address"`
	Notes        string `json:"notes"`
}

type UpdateBuyerRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=150"`
	ContactName  string `json:"contact_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"omitempty,email"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email"`
	Address      string `json:"address"`
	Notes        string `json:"notes"`
}

type BuyerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	Address   string    `json:"address,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BuyerSummaryResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Phone string    `json:"phone,omitempty"`
}

type DeletedBuyerResponse struct {
	ID      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}

func NewBuyerResponse(buyer models.Buyer) BuyerResponse {
	return BuyerResponse{
		ID:        buyer.ID,
		Name:      buyer.Name,
		Phone:     buyer.ContactPhone,
		Email:     buyer.ContactEmail,
		Address:   buyer.Address,
		Notes:     buyer.Notes,
		CreatedAt: buyer.CreatedAt,
		UpdatedAt: buyer.UpdatedAt,
	}
}

func NewBuyerResponses(buyers []models.Buyer) []BuyerResponse {
	responses := make([]BuyerResponse, len(buyers))
	for i, buyer := range buyers {
		responses[i] = NewBuyerResponse(buyer)
	}
	return responses
}

func NewBuyerSummaryResponse(buyer models.Buyer) *BuyerSummaryResponse {
	return &BuyerSummaryResponse{
		ID:    buyer.ID,
		Name:  buyer.Name,
		Phone: buyer.ContactPhone,
	}
}
