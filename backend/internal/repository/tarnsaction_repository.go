package repository

import (
	"context"
	"seagame/ticket/backend/internal/models/transaction"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *transaction.Transaction) error
	DailyIncome(ctx context.Context, date time.Time) ([]transaction.Transaction, error)
	WeeklyIncome(ctx context.Context, from, to time.Time) ([]transaction.Transaction, error)
	MonthlyIncome(ctx context.Context, year, month int) ([]transaction.Transaction, error)
}

type transactionRepo struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db}
}

func (r *transactionRepo) Create(ctx context.Context, tx *transaction.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepo) DailyIncome(ctx context.Context, date time.Time) ([]transaction.Transaction, error) {
	var txs []transaction.Transaction
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	return txs, r.db.WithContext(ctx).
		Where("processed_at >= ? AND processed_at < ?", start, end).
		Find(&txs).Error
}

func (r *transactionRepo) WeeklyIncome(ctx context.Context, from, to time.Time) ([]transaction.Transaction, error) {
	var txs []transaction.Transaction
	return txs, r.db.WithContext(ctx).
		Where("processed_at >= ? AND processed_at <= ?", from, to).
		Find(&txs).Error
}

func (r *transactionRepo) MonthlyIncome(ctx context.Context, year, month int) ([]transaction.Transaction, error) {
	var txs []transaction.Transaction
	return txs, r.db.WithContext(ctx).
		Where(
			"EXTRACT(YEAR FROM processed_at) = ? AND EXTRACT(MONTH FROM processed_at) = ?",
			year, month,
		).
		Find(&txs).Error
}
