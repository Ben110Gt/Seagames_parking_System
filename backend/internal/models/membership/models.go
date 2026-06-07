package membership

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipCard struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	CardCode         string    `gorm:"uniqueIndex;not null"`
	PlateNumber      string    `gorm:"not null;index"`
	RegistrationDate time.Time `gorm:"not null"`
	ExpirationDate   time.Time `gorm:"not null"`
	FeeAmount        int64     `gorm:"default:60000"`
	IssuedBy         uuid.UUID `gorm:"type:uuid;not null"` // FK -> User.ID
}

func (m *MembershipCard) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

func (m *MembershipCard) IsActive() bool {
	return time.Now().Before(m.ExpirationDate)
}
