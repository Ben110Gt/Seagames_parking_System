package repository

import (
	"context"
	"errors"
	"fmt"
	"seagame/ticket/backend/internal/models/ticket"
	"time"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	Create(ctx context.Context, m *ticket.Membership) error
	FindActiveByPlateNumber(ctx context.Context, plateNumber string) (*ticket.Membership, error)
	Update(ctx context.Context, m *ticket.Membership) error
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) Create(ctx context.Context, m *ticket.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// FindActiveByPlateNumber ค้นหา membership ที่ยังไม่หมดอายุ
func (r *membershipRepository) FindActiveByPlateNumber(ctx context.Context, plateNumber string) (*ticket.Membership, error) {
	var m ticket.Membership
	err := r.db.WithContext(ctx).
		Where("plate_number = ? AND status = ? AND end_date > ?", plateNumber, "active", time.Now()).
		Order("end_date DESC").
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no active membership for plate: %s", plateNumber)
		}
		return nil, err
	}
	return &m, nil
}

func (r *membershipRepository) Update(ctx context.Context, m *ticket.Membership) error {
	return r.db.WithContext(ctx).Save(m).Error
}
