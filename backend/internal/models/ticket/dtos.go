package ticket

import (
	"time"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	// CustomerRole string    `json:"customer_role" validate:"required,oneof=Customer Membership"`
	IssuedBy uuid.UUID `json:"-"`
}

type CreateTicketResponse struct {
	Ticket *Ticket `json:"ticket"`
	QRCode string  `json:"qr_code"`
}

type CheckTicketRequest struct {
	TicketCode string    `json:"ticket_code" validate:"required"`
	CheckedBy  uuid.UUID `json:"-"`
}

type CheckTicketResponse struct {
	Ticket   *Ticket           `json:"ticket,omitempty"`
	Checkout *CheckoutResponse `json:"checkout,omitempty"`
	Message  string            `json:"message"`
}

type TicketResponse struct {
	ID           uuid.UUID  `json:"id"`
	TicketCode   string     `json:"ticket_code"`
	PlateNumber  string     `json:"plate_number"`
	CustomerRole string     `json:"customer_role"`
	CheckinTime  time.Time  `json:"checkin_time"`
	CheckoutTime *time.Time `json:"checkout_time,omitempty"`
	Status       string     `json:"status"`
	FineAmount   int64      `json:"fine_amount"`
	BarcodeURL   string     `json:"barcode_url"`
}

type CheckoutResponse struct {
	TicketCode   string    `json:"ticket_code"`
	PlateNumber  string    `json:"plate_number"`
	CustomerRole string    `json:"customer_role"`
	CheckinTime  time.Time `json:"checkin_time"`
	CheckoutTime time.Time `json:"checkout_time"`
	DaysParked   int       `json:"days_parked"`
	FineAmount   int64     `json:"fine_amount"`
	Message      string    `json:"message"`
}

type SearchTicketRequest struct {
	PlateNumber string `json:"plate_number"`
	Query       string `query:"q"`
	Status      string `query:"status"`
}

type SearchTicketResponse struct {
	Tickets []Ticket `json:"tickets"`
	Count   int      `json:"count"`
}

type IncomeRequest struct {
	Period string `json:"period"`
}

type IncomeDetail struct {
	Date   string  `json:"date"`
	Income float64 `json:"income"`
	Count  int     `json:"count"`
}

type IncomeResponse struct {
	Period      string         `json:"period"`
	TotalIncome float64        `json:"total_income"`
	TotalCount  int            `json:"total_count"`
	Details     []IncomeDetail `json:"details"`
}
