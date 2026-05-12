package ticket

type CreateTicketRequest struct {
	PlateNumber string `json:"plate_number" validate:"required,min=2,max=20"`
}

type CreateTicketResponse struct {
	Ticket *Ticket `json:"ticket"`
	QRCode string  `json:"qr_code"`
}

type CheckTicketRequest struct {
	TicketCode string `json:"ticket_code" validate:"required"`
}

type CheckTicketResponse struct {
	Ticket  *Ticket `json:"ticket"`
	Message string  `json:"message"`
}

// ค้นหาตั๋วจากทะเบียนรถ
type SearchTicketRequest struct {
	PlateNumber string `json:"plate_number" query:"plate_number" validate:"required"`
}

type SearchTicketResponse struct {
	Tickets []Ticket `json:"tickets"`
	Count   int      `json:"count"`
}

// รายได้ dashboard
type IncomeRequest struct {
	Period string `json:"period" query:"period"` // daily, weekly, monthly
}

type IncomeResponse struct {
	Period      string        `json:"period"`
	TotalIncome float64       `json:"total_income"`
	TotalCount  int           `json:"total_count"`
	Details     []IncomeDetail `json:"details"`
}

type IncomeDetail struct {
	Date   string  `json:"date"`
	Income float64 `json:"income"`
	Count  int     `json:"count"`
}
