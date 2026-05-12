package service

import (
	"context"
	"fmt"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/repository"
	"time"
)

const MembershipFee float64 = 60000 // 60,000 KIP ต่อเดือน

type MembershipService interface {
	CreateMembership(ctx context.Context, req *ticket.CreateMembershipRequest) (*ticket.CreateMembershipResponse, error)
	CheckMembership(ctx context.Context, req *ticket.CheckMembershipRequest) (*ticket.CheckMembershipResponse, error)
}

type membershipService struct {
	repo repository.MembershipRepository
}

func NewMembershipService(repo repository.MembershipRepository) MembershipService {
	return &membershipService{repo: repo}
}

// CreateMembership สร้าง/ต่ออายุ membership รายเดือน
func (s *membershipService) CreateMembership(ctx context.Context, req *ticket.CreateMembershipRequest) (*ticket.CreateMembershipResponse, error) {
	now := time.Now()
	startDate := now

	// ถ้ามี membership เดิมที่ยังใช้ได้ ให้ต่อจากวันหมดอายุเดิม
	existing, err := s.repo.FindActiveByPlateNumber(ctx, req.PlateNumber)
	if err == nil && existing != nil && existing.IsActive() {
		startDate = existing.EndDate
		// ปิด membership เก่า
		existing.Status = "renewed"
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update old membership: %w", err)
		}
	}

	endDate := startDate.AddDate(0, 1, 0) // +1 เดือน

	m := &ticket.Membership{
		PlateNumber: req.PlateNumber,
		StartDate:   startDate,
		EndDate:     endDate,
		Fee:         MembershipFee,
		Status:      "active",
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return &ticket.CreateMembershipResponse{
		Membership: m,
		Message:    fmt.Sprintf("Membership created for %s. Valid until %s", req.PlateNumber, endDate.Format("2006-01-02")),
	}, nil
}

// CheckMembership ตรวจสอบสถานะ membership
func (s *membershipService) CheckMembership(ctx context.Context, req *ticket.CheckMembershipRequest) (*ticket.CheckMembershipResponse, error) {
	m, err := s.repo.FindActiveByPlateNumber(ctx, req.PlateNumber)
	if err != nil {
		return &ticket.CheckMembershipResponse{
			IsMember:   false,
			Membership: nil,
		}, nil
	}
	return &ticket.CheckMembershipResponse{
		IsMember:   m.IsActive(),
		Membership: m,
	}, nil
}
