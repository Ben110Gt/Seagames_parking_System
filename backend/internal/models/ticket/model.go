package ticket

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TicketCode   string         `gorm:"uniqueIndex;not null"`
	PlateNumber  string         `gorm:"not null;index"`
	CustomerRole string         `gorm:"not null;default:Customer"` // Customer | Membership
	CheckinTime  time.Time      `gorm:"not null"`
	CheckoutTime *time.Time     `gorm:"index"`
	Status       string         `gorm:"not null;default:in;index"` // in | out
	FineAmount   int64          `gorm:"default:0"`                 // kip
	IssuedBy     uuid.UUID      `gorm:"type:uuid;not null"`        // FK -> User.ID
	CheckedBy    *uuid.UUID     `gorm:"type:uuid"`                 // FK -> User.ID
	DeletedAt    gorm.DeletedAt `gorm:"index"`                     // soft delete
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
