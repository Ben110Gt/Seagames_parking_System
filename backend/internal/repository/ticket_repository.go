package repository

import (
	"context"
	"errors"
	"fmt"
	"seagame/ticket/backend/internal/models/ticket"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ctx context.Context, t *ticket.Ticket) error
	FindByTicketCode(ctx context.Context, ticketCode string) (ticket.Ticket, error)
	FindActiveByPlateNumber(ctx context.Context, plateNumber string) (ticket.Ticket, error)
	Update(ctx context.Context, t *ticket.Ticket) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SearchByPlateNumber(ctx context.Context, plateNumber string) ([]ticket.Ticket, error)
	GetIncomeByDateRange(ctx context.Context, start, end time.Time) ([]ticket.Ticket, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, t *ticket.Ticket) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *ticketRepository) FindByTicketCode(ctx context.Context, ticketCode string) (ticket.Ticket, error) {
	var tk ticket.Ticket
	err := r.db.WithContext(ctx).Where("ticket_code = ?", ticketCode).First(&tk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tk, fmt.Errorf("ticket not found: %s", ticketCode)
		}
		return tk, err
	}
	return tk, nil
}

func (r *ticketRepository) FindActiveByPlateNumber(ctx context.Context, plateNumber string) (ticket.Ticket, error) {
	var tk ticket.Ticket
	err := r.db.WithContext(ctx).
		Where("plate_number = ? AND status = ?", plateNumber, "In").
		Order("created_at DESC").
		First(&tk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tk, fmt.Errorf("no active ticket for plate: %s", plateNumber)
		}
		return tk, err
	}
	return tk, nil
}

func (r *ticketRepository) Update(ctx context.Context, t *ticket.Ticket) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *ticketRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ticket.Ticket{}, "id = ?", id).Error
}

// SearchByPlateNumber ค้นหาตั๋วจากทะเบียนรถ (กรณีลูกค้าทำตั๋วหาย)
func (r *ticketRepository) SearchByPlateNumber(ctx context.Context, plateNumber string) ([]ticket.Ticket, error) {
	var tickets []ticket.Ticket
	err := r.db.WithContext(ctx).
		Where("plate_number ILIKE ?", "%"+plateNumber+"%").
		Order("created_at DESC").
		Limit(50).
		Find(&tickets).Error
	return tickets, err
}

// GetIncomeByDateRange ดึงตั๋วที่ check-out แล้วในช่วงเวลา
func (r *ticketRepository) GetIncomeByDateRange(ctx context.Context, start, end time.Time) ([]ticket.Ticket, error) {
	var tickets []ticket.Ticket
	err := r.db.WithContext(ctx).
		Unscoped(). // รวมที่ soft-deleted ด้วย (ตั๋วที่ check-out แล้ว)
		Where("status = ? AND check_out BETWEEN ? AND ?", "Out", start, end).
		Find(&tickets).Error
	return tickets, err
}
