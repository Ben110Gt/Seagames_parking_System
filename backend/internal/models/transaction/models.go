package transaction

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction records a completed payment event (fine or membership)
type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReferenceCode string    `gorm:"not null;index"` // TicketCode or CardCode
	Type          string    `gorm:"not null"`       // fine | membership
	Amount        int64     `gorm:"not null"`
	ProcessedAt   time.Time `gorm:"not null"`
	ProcessedBy   uuid.UUID `gorm:"type:uuid;not null"` // FK -> User.ID
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
