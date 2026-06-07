package service

import (
	"context"
	"seagame/ticket/backend/internal/models/transaction"
	"seagame/ticket/backend/internal/repository"
	"time"
)

type IncomeService interface {
	DailyReport(ctx context.Context, date time.Time) (*transaction.IncomeReport, error)
	WeeklyReport(ctx context.Context, weekStart time.Time) (*transaction.IncomeReport, error)
	MonthlyReport(ctx context.Context, year, month int) (*transaction.IncomeReport, error)
}

type incomeService struct {
	transactionRepo repository.TransactionRepository
}

func NewIncomeService(transactionRepo repository.TransactionRepository) IncomeService {
	return &incomeService{
		transactionRepo: transactionRepo,
	}
}

func (s *incomeService) DailyReport(ctx context.Context, date time.Time) (*transaction.IncomeReport, error) {
	txs, err := s.transactionRepo.DailyIncome(ctx, date)
	if err != nil {
		return nil, err
	}
	return buildReport("daily", txs), nil
}

func (s *incomeService) WeeklyReport(ctx context.Context, weekStart time.Time) (*transaction.IncomeReport, error) {
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	txs, err := s.transactionRepo.WeeklyIncome(ctx, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	return buildReport("weekly", txs), nil
}

func (s *incomeService) MonthlyReport(ctx context.Context, year, month int) (*transaction.IncomeReport, error) {
	txs, err := s.transactionRepo.MonthlyIncome(ctx, year, month)
	if err != nil {
		return nil, err
	}
	return buildReport("monthly", txs), nil
}

func buildReport(period string, txs []transaction.Transaction) *transaction.IncomeReport {
	r := &transaction.IncomeReport{Period: period}
	dailyMap := map[string]*transaction.IncomeLine{}

	for _, tx := range txs {
		r.TotalIncome += tx.Amount
		if tx.Type == "fine" {
			r.TotalFines += tx.Amount
		}
		if tx.Type == "membership" {
			r.MembershipRevenue += tx.Amount
		}
		r.Transactions++

		day := tx.ProcessedAt.Format("2006-01-02")
		if _, ok := dailyMap[day]; !ok {
			dailyMap[day] = &transaction.IncomeLine{Date: day}
		}
		dailyMap[day].Income += tx.Amount
		dailyMap[day].Transactions++
		if tx.Type == "fine" {
			dailyMap[day].Fines += tx.Amount
		}
	}

	for _, line := range dailyMap {
		r.Breakdown = append(r.Breakdown, *line)
	}
	return r
}
