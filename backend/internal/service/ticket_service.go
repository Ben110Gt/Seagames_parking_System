package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/repository"
	util "seagame/ticket/backend/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketService interface {
	CreateTicket(ctx context.Context, req *ticket.CreateTicketRequest) (*ticket.CreateTicketResponse, error)
	CheckTicket(ctx context.Context, req *ticket.CheckTicketRequest) (*ticket.CheckTicketResponse, error)
	SearchTicket(ctx context.Context, req *ticket.SearchTicketRequest) (*ticket.SearchTicketResponse, error)
	GetIncome(ctx context.Context, req *ticket.IncomeRequest) (*ticket.IncomeResponse, error)
}

type ticketService struct {
	repo            repository.TicketRepository
	memberRepo      repository.MembershipRepository
	transactionRepo repository.TransactionRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, memberRepo repository.MembershipRepository, transactionRepo repository.TransactionRepository) TicketService {
	return &ticketService{
		repo:            ticketRepo,
		memberRepo:      memberRepo,
		transactionRepo: transactionRepo,
	}
}

const (
	BaseFee    float64 = 2000
	PenaltyFee float64 = 2000
)

func (s *ticketService) CreateTicket(ctx context.Context, req *ticket.CreateTicketRequest) (*ticket.CreateTicketResponse, error) {
	ticketCode := util.GenerateTicketCode()
	now := time.Now()

	tk := &ticket.Ticket{
		TicketCode:   ticketCode,
		PlateNumber:  req.PlateNumber,
		CustomerRole: "Customer",
		CheckinTime:  now,
		Status:       "in",
		IssuedBy:     req.IssuedBy,
	}

	if err := s.repo.Create(ctx, tk); err != nil {
		return nil, err
	}

	qrCode, err := util.GenerateQRCodeBase64(ticketCode)
	if err != nil {
		return nil, err
	}

	return &ticket.CreateTicketResponse{
		Ticket: tk,
		QRCode: qrCode,
	}, nil
}

func (s *ticketService) checkoutMembership(ctx context.Context, cardCode string) (*ticket.CheckTicketResponse, error) {
	// if s.memberRepo == nil {
	// 	return nil, errors.New("ticket or membership card not found")
	// }

	card, err := s.memberRepo.FindByCode(ctx, cardCode)
	if err != nil {
		return nil, errors.New("ticket or membership card not found")
	}

	now := time.Now()
	if !card.IsActive() {
		return nil, errors.New("membership expired — please purchase a day ticket")
	}

	checkout := &ticket.CheckoutResponse{
		TicketCode:   card.CardCode,
		PlateNumber:  card.PlateNumber,
		CustomerRole: "Membership",
		CheckinTime:  card.RegistrationDate,
		CheckoutTime: now,
		DaysParked:   0,
		FineAmount:   0,
		Message:      "Membership valid. Access granted. No charge.",
	}

	return &ticket.CheckTicketResponse{
		Checkout: checkout,
		Message:  checkout.Message,
	}, nil
}

func (s *ticketService) CheckTicket(ctx context.Context, req *ticket.CheckTicketRequest) (*ticket.CheckTicketResponse, error) {
	tk, err := s.repo.FindByCode(ctx, req.TicketCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if archived, archErr := s.repo.FindByCodeUnscoped(ctx, req.TicketCode); archErr == nil {
				return nil, fmt.Errorf("ticket %s has already been checked out", archived.TicketCode)
			}
			if strings.HasPrefix(req.TicketCode, "MBR-") {
				return s.checkoutMembership(ctx, req.TicketCode)
			}
			return nil, fmt.Errorf("ticket %s not found", req.TicketCode)
		}
		return nil, err
	}

	if tk.Status == "out" {
		return nil, fmt.Errorf("ticket %s has already been checked out", req.TicketCode)
	}

	now := time.Now()
	fineAmount := int64(0)

	isMember := false
	if s.memberRepo != nil {
		member, err := s.memberRepo.FindActiveByPlateNumber(ctx, tk.PlateNumber)
		if err == nil && member != nil && member.IsActive() {
			isMember = true
		}
	}

	if !isMember {
		fineAmount = int64(calculatePenalty(&tk.CheckinTime, &now))
	}

	tk.CheckoutTime = &now
	tk.FineAmount = fineAmount
	tk.Status = "out"
	if req.CheckedBy != uuid.Nil {
		tk.CheckedBy = &req.CheckedBy
	}

	if err := s.repo.UpdateCheckout(ctx, tk); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	if err := s.repo.SoftDelete(ctx, tk); err != nil {
		return nil, fmt.Errorf("failed to archive ticket: %w", err)
	}

	totalFee := BaseFee + float64(fineAmount)
	message := fmt.Sprintf("Check-out successful. Total fee: %.0f KIP", totalFee)
	if !isMember && fineAmount > 0 {
		days := int(fineAmount) / int(PenaltyFee)
		message = fmt.Sprintf(
			"Check-out successful. Late %d day(s). Penalty: %d KIP. Total: %.0f KIP",
			days, fineAmount, totalFee,
		)
	} else if isMember {
		message = fmt.Sprintf("Check-out successful. Member — no penalty. Fee: %.0f KIP", BaseFee)
	}

	return &ticket.CheckTicketResponse{
		Ticket:  tk,
		Message: message,
	}, nil
}

func (s *ticketService) SearchTicket(ctx context.Context, req *ticket.SearchTicketRequest) (*ticket.SearchTicketResponse, error) {
	var tickets []ticket.Ticket
	var err error

	if req.PlateNumber != "" {
		tickets, err = s.repo.FindByPlate(ctx, req.PlateNumber)
	} else {
		tickets, err = s.repo.SearchTickets(ctx, req.Query, req.Status)
	}
	if err != nil {
		return nil, err
	}

	return &ticket.SearchTicketResponse{
		Tickets: tickets,
		Count:   len(tickets),
	}, nil
}

func (s *ticketService) GetIncome(ctx context.Context, req *ticket.IncomeRequest) (*ticket.IncomeResponse, error) {
	now := time.Now()
	var start time.Time

	switch req.Period {
	case "weekly":
		start = now.AddDate(0, 0, -7)
	case "monthly":
		start = now.AddDate(0, -1, 0)
	default:
		start = truncateToDay(now)
		req.Period = "daily"
	}

	tickets, err := s.repo.GetIncomeByDateRange(ctx, start, now)
	if err != nil {
		return nil, err
	}

	var totalIncome float64
	dayMap := make(map[string]*ticket.IncomeDetail)

	for _, t := range tickets {
		income := BaseFee + float64(t.FineAmount)
		totalIncome += income

		dateKey := t.CheckinTime.Format("2006-01-02")
		if t.CheckoutTime != nil {
			dateKey = t.CheckoutTime.Format("2006-01-02")
		}

		if detail, ok := dayMap[dateKey]; ok {
			detail.Income += income
			detail.Count++
		} else {
			dayMap[dateKey] = &ticket.IncomeDetail{
				Date:   dateKey,
				Income: income,
				Count:  1,
			}
		}
	}

	details := make([]ticket.IncomeDetail, 0, len(dayMap))
	for _, d := range dayMap {
		details = append(details, *d)
	}

	return &ticket.IncomeResponse{
		Period:      req.Period,
		TotalIncome: totalIncome,
		TotalCount:  len(tickets),
		Details:     details,
	}, nil
}

func calculatePenalty(checkIn, checkOut *time.Time) float64 {
	if checkIn == nil || checkOut == nil {
		return 0
	}
	inDate := truncateToDay(*checkIn)
	outDate := truncateToDay(*checkOut)

	days := int(math.Floor(outDate.Sub(inDate).Hours() / 24))
	if days <= 0 {
		return 0
	}
	return float64(days) * PenaltyFee
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
