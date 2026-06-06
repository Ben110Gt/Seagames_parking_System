package repository

import (
	"context"
	"errors"
	"fmt"
	"seagame/ticket/backend/internal/models/membership"
	"time"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	Create(ctx context.Context, m *membership.MembershipCard) error
	FindByCode(ctx context.Context, code string) (*membership.MembershipCard, error)
	FindByPlate(ctx context.Context, plate string) (*membership.MembershipCard, error)
	FindActiveByPlateNumber(ctx context.Context, plateNumber string) (*membership.MembershipCard, error)
	FindAll(ctx context.Context) ([]membership.MembershipCard, error)
	FindActive(ctx context.Context) ([]membership.MembershipCard, error)
}

type membershipRepo struct{ db *gorm.DB }

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepo{db}
}

func (r *membershipRepo) Create(ctx context.Context, m *membership.MembershipCard) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *membershipRepo) FindByCode(ctx context.Context, code string) (*membership.MembershipCard, error) {
	var m membership.MembershipCard
	return &m, r.db.WithContext(ctx).Where("card_code = ?", code).First(&m).Error
}

func (r *membershipRepo) FindByPlate(ctx context.Context, plate string) (*membership.MembershipCard, error) {
	var m membership.MembershipCard
	return &m, r.db.WithContext(ctx).
		Where("plate_number ILIKE ? AND expiration_date >= ?", plate, time.Now()).
		First(&m).Error
}

func (r *membershipRepo) FindActiveByPlateNumber(ctx context.Context, plateNumber string) (*membership.MembershipCard, error) {
	var m membership.MembershipCard
	err := r.db.WithContext(ctx).
		Where("plate_number = ? AND expiration_date > ?", plateNumber, time.Now()).
		Order("expiration_date DESC").
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no active membership for plate: %s", plateNumber)
		}
		return nil, err
	}
	return &m, nil
}

func (r *membershipRepo) FindAll(ctx context.Context) ([]membership.MembershipCard, error) {
	var ms []membership.MembershipCard
	return ms, r.db.WithContext(ctx).Order("registration_date DESC").Find(&ms).Error
}

func (r *membershipRepo) FindActive(ctx context.Context) ([]membership.MembershipCard, error) {
	var ms []membership.MembershipCard
	return ms, r.db.WithContext(ctx).
		Where("expiration_date >= ?", time.Now()).
		Find(&ms).Error
}
