package service

import (
	"context"
	"fmt"
	"math"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/repository"
	util "seagame/ticket/backend/utils"
	"time"
)

type TicketService interface {
	CreateTicket(ctx context.Context, req *ticket.CreateTicketRequest) (*ticket.CreateTicketResponse, error)
	CheckTicket(ctx context.Context, req *ticket.CheckTicketRequest) (*ticket.CheckTicketResponse, error)
	SearchTicket(ctx context.Context, req *ticket.SearchTicketRequest) (*ticket.SearchTicketResponse, error)
	GetIncome(ctx context.Context, req *ticket.IncomeRequest) (*ticket.IncomeResponse, error)
}

type ticketService struct {
	repo       repository.TicketRepository
	memberRepo repository.MembershipRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, memberRepo repository.MembershipRepository) TicketService {
	return &ticketService{
		repo:       ticketRepo,
		memberRepo: memberRepo,
	}
}

const (
	BaseFee    float64 = 2000
	PenaltyFee float64 = 2000 // ค่าปรับต่อวัน
)

// CreateTicket ออกตั๋วเมื่อรถเข้า (Check-in)
func (s *ticketService) CreateTicket(ctx context.Context, req *ticket.CreateTicketRequest) (*ticket.CreateTicketResponse, error) {
	ticketCode := fmt.Sprintf("SG-%d", time.Now().UnixNano())
	now := time.Now()

	tk := &ticket.Ticket{
		TicketCode:  ticketCode,
		PlateNumber: req.PlateNumber,
		CheckIn:     &now,
		TotalFee:    BaseFee,
		Status:      "In", // รถเข้ามาแล้ว
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

// CheckTicket สแกนตั๋วเมื่อรถออก (Check-out)
func (s *ticketService) CheckTicket(ctx context.Context, req *ticket.CheckTicketRequest) (*ticket.CheckTicketResponse, error) {
	tk, err := s.repo.FindByTicketCode(ctx, req.TicketCode)
	if err != nil {
		return nil, err
	}

	if tk.Status == "Out" {
		return nil, fmt.Errorf("ticket %s has already been checked out", req.TicketCode)
	}

	now := time.Now()
	totalFee := BaseFee

	// ตรวจสอบ membership — ถ้าเป็นสมาชิกไม่คิดค่าปรับ
	isMember := false
	if s.memberRepo != nil {
		member, err := s.memberRepo.FindActiveByPlateNumber(ctx, tk.PlateNumber)
		if err == nil && member != nil && member.IsActive() {
			isMember = true
		}
	}

	if !isMember {
		penalty := calculatePenalty(tk.CheckIn, &now)
		totalFee += penalty
	}

	tk.CheckOut = &now
	tk.TotalFee = totalFee
	tk.Status = "Out"

	if err := s.repo.Update(ctx, &tk); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	// Soft delete (archive) ตั๋วที่ check-out แล้ว
	if err := s.repo.SoftDelete(ctx, tk.ID); err != nil {
		return nil, fmt.Errorf("failed to archive ticket: %w", err)
	}

	message := fmt.Sprintf("Check-out successful. Total fee: %.0f KIP", totalFee)
	if !isMember {
		penalty := calculatePenalty(tk.CheckIn, &now)
		if penalty > 0 {
			days := int(penalty / PenaltyFee)
			message = fmt.Sprintf(
				"Check-out successful. Late %d day(s). Penalty: %.0f KIP. Total: %.0f KIP",
				days, penalty, totalFee,
			)
		}
	} else {
		message = fmt.Sprintf("Check-out successful. Member — no penalty. Fee: %.0f KIP", BaseFee)
	}

	return &ticket.CheckTicketResponse{
		Ticket:  &tk,
		Message: message,
	}, nil
}

// SearchTicket ค้นหาตั๋วจากทะเบียนรถ (กรณีลูกค้าทำตั๋วหาย)
func (s *ticketService) SearchTicket(ctx context.Context, req *ticket.SearchTicketRequest) (*ticket.SearchTicketResponse, error) {
	tickets, err := s.repo.SearchByPlateNumber(ctx, req.PlateNumber)
	if err != nil {
		return nil, err
	}
	return &ticket.SearchTicketResponse{
		Tickets: tickets,
		Count:   len(tickets),
	}, nil
}

// GetIncome ดูรายได้ (daily / weekly / monthly)
func (s *ticketService) GetIncome(ctx context.Context, req *ticket.IncomeRequest) (*ticket.IncomeResponse, error) {
	now := time.Now()
	var start time.Time

	switch req.Period {
	case "weekly":
		start = now.AddDate(0, 0, -7)
	case "monthly":
		start = now.AddDate(0, -1, 0)
	default: // daily
		start = truncateToDay(now)
		req.Period = "daily"
	}

	tickets, err := s.repo.GetIncomeByDateRange(ctx, start, now)
	if err != nil {
		return nil, err
	}

	// คำนวณรายได้รวม
	var totalIncome float64
	dayMap := make(map[string]*ticket.IncomeDetail)

	for _, t := range tickets {
		totalIncome += t.TotalFee
		dateKey := ""
		if t.CheckOut != nil {
			dateKey = t.CheckOut.Format("2006-01-02")
		} else {
			dateKey = t.CreatedAt.Format("2006-01-02")
		}

		if detail, ok := dayMap[dateKey]; ok {
			detail.Income += t.TotalFee
			detail.Count++
		} else {
			dayMap[dateKey] = &ticket.IncomeDetail{
				Date:   dateKey,
				Income: t.TotalFee,
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

// calculatePenalty คำนวณค่าปรับตามจำนวนวันที่เกิน
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
