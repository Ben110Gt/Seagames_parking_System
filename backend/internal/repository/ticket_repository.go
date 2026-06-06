package repository

import (
	"context"
	"seagame/ticket/backend/internal/models/ticket"
	"time"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ctx context.Context, t *ticket.Ticket) error
	FindByCode(ctx context.Context, code string) (*ticket.Ticket, error)
	FindByCodeUnscoped(ctx context.Context, code string) (*ticket.Ticket, error)
	FindByPlate(ctx context.Context, plate string) ([]ticket.Ticket, error)
	FindActive(ctx context.Context) ([]ticket.Ticket, error)
	SearchTickets(ctx context.Context, query, status string) ([]ticket.Ticket, error)
	GetIncomeByDateRange(ctx context.Context, start, end time.Time) ([]ticket.Ticket, error)
	UpdateCheckout(ctx context.Context, t *ticket.Ticket) error
	SoftDelete(ctx context.Context, t *ticket.Ticket) error
}

type ticketRepo struct{ db *gorm.DB }

func NewTicketRepository(db *gorm.DB) TicketRepository { return &ticketRepo{db} }

func (r *ticketRepo) Create(ctx context.Context, t *ticket.Ticket) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *ticketRepo) FindByCode(ctx context.Context, code string) (*ticket.Ticket, error) {
	var t ticket.Ticket
	return &t, r.db.WithContext(ctx).Where("ticket_code = ?", code).First(&t).Error
}

func (r *ticketRepo) FindByCodeUnscoped(ctx context.Context, code string) (*ticket.Ticket, error) {
	var t ticket.Ticket
	return &t, r.db.WithContext(ctx).Unscoped().Where("ticket_code = ?", code).First(&t).Error
}

func (r *ticketRepo) FindByPlate(ctx context.Context, plate string) ([]ticket.Ticket, error) {
	var ts []ticket.Ticket
	return ts, r.db.WithContext(ctx).Where("plate_number ILIKE ?", "%"+plate+"%").Order("checkin_time DESC").Find(&ts).Error
}

func (r *ticketRepo) FindActive(ctx context.Context) ([]ticket.Ticket, error) {
	var ts []ticket.Ticket
	return ts, r.db.WithContext(ctx).Where("status = ?", "in").Order("checkin_time ASC").Find(&ts).Error
}

func (r *ticketRepo) SearchTickets(ctx context.Context, query, status string) ([]ticket.Ticket, error) {
	var ts []ticket.Ticket
	q := r.db.WithContext(ctx).Where("plate_number ILIKE ? OR ticket_code ILIKE ?", "%"+query+"%", "%"+query+"%")
	if status != "" && status != "all" {
		q = q.Where("status = ?", status)
	}
	return ts, q.Order("checkin_time DESC").Limit(50).Find(&ts).Error
}

func (r *ticketRepo) GetIncomeByDateRange(ctx context.Context, start, end time.Time) ([]ticket.Ticket, error) {
	var ts []ticket.Ticket
	return ts, r.db.WithContext(ctx).Unscoped().
		Where("status = ? AND checkout_time >= ? AND checkout_time <= ?", "out", start, end).
		Find(&ts).Error
}

func (r *ticketRepo) UpdateCheckout(ctx context.Context, t *ticket.Ticket) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *ticketRepo) SoftDelete(ctx context.Context, t *ticket.Ticket) error {
	return r.db.WithContext(ctx).Delete(t).Error
}
