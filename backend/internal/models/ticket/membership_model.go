package ticket

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Membership ระบบสมาชิกรายเดือน 60,000 KIP
type Membership struct {
	ID          uuid.UUID      `json:"id"           gorm:"type:uuid;primaryKey"`
	PlateNumber string         `json:"plate_number" gorm:"not null;index"`
	StartDate   time.Time      `json:"start_date"   gorm:"not null"`
	EndDate     time.Time      `json:"end_date"     gorm:"not null"`
	Fee         float64        `json:"fee"          gorm:"default:60000"`
	Status      string         `json:"status"       gorm:"type:varchar(20);default:'active'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"            gorm:"index"`
}

func (m *Membership) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// IsActive ตรวจสอบว่า membership ยังใช้งานอยู่ไหม
func (m *Membership) IsActive() bool {
	return m.Status == "active" && time.Now().Before(m.EndDate)
}

// DTOs

type CreateMembershipRequest struct {
	PlateNumber string `json:"plate_number" validate:"required"`
}

type CreateMembershipResponse struct {
	Membership *Membership `json:"membership"`
	Message    string      `json:"message"`
}

type CheckMembershipRequest struct {
	PlateNumber string `json:"plate_number" validate:"required"`
}

type CheckMembershipResponse struct {
	IsMember   bool        `json:"is_member"`
	Membership *Membership `json:"membership,omitempty"`
}
