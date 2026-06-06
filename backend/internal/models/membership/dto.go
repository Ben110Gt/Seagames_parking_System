package membership

import (
	"time"

	"github.com/google/uuid"
)

type CreateMembershipRequest struct {
	PlateNumber string `json:"plate_number" validate:"required"`
}

type MembershipResponse struct {
	ID               uuid.UUID `json:"id"`
	CardCode         string    `json:"card_code"`
	PlateNumber      string    `json:"plate_number"`
	RegistrationDate time.Time `json:"registration_date"`
	ExpirationDate   time.Time `json:"expiration_date"`
	IsActive         bool      `json:"is_active"`
	BarcodeURL       string    `json:"barcode_url"`
}
