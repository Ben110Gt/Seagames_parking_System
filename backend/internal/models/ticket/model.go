package ticket

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID          uuid.UUID      `json:"id"           gorm:"type:uuid;primaryKey"`
	TicketCode  string         `json:"ticket_code"  gorm:"uniqueIndex;not null"`
	PlateNumber string         `json:"plate_number" gorm:"not null;index"`
	CheckIn     *time.Time     `json:"check_in"`
	CheckOut    *time.Time     `json:"check_out"`
	TotalFee    float64        `json:"total_fee"    gorm:"default:2000"`
	Status      string         `json:"status"       gorm:"type:varchar(10);default:'In'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"            gorm:"index"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
