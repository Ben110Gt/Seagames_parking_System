package service

import (
	"context"
	"fmt"
	"os"
	"seagame/ticket/backend/internal/models/membership"
	"seagame/ticket/backend/internal/models/transaction"
	"seagame/ticket/backend/internal/repository"
	util "seagame/ticket/backend/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type MembershipService interface {
	CreateMembership(ctx context.Context, req *membership.CreateMembershipRequest, issuedBy uuid.UUID) (*membership.MembershipResponse, error)
	GetMembershipByCode(ctx context.Context, code string) (*membership.MembershipResponse, error)
	GetAllMemberships(ctx context.Context) ([]membership.MembershipResponse, error)
	GetActiveMemberships(ctx context.Context) ([]membership.MembershipResponse, error)
}

type membershipService struct {
	membershipRepo  repository.MembershipRepository
	transactionRepo repository.TransactionRepository
}

func NewMembershipService(membershipRepo repository.MembershipRepository, transactionRepo repository.TransactionRepository) MembershipService {
	return &membershipService{
		membershipRepo:  membershipRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *membershipService) CreateMembership(ctx context.Context, req *membership.CreateMembershipRequest, issuedBy uuid.UUID) (*membership.MembershipResponse, error) {
	code := util.GenerateMemberCode()
	now := time.Now()

	fee, _ := strconv.ParseInt(os.Getenv("MEMBERSHIP_FEE_KIP"), 10, 64)
	if fee == 0 {
		fee = 60000
	}

	loc, _ := time.LoadLocation("Asia/Vientiane")
	y, m, _ := now.In(loc).Date()
	exp := time.Date(y, m+1, 1, 0, 0, 0, 0, loc).Add(-time.Second)

	card := &membership.MembershipCard{
		CardCode:         code,
		PlateNumber:      req.PlateNumber,
		RegistrationDate: now,
		ExpirationDate:   exp,
		FeeAmount:        fee,
		IssuedBy:         issuedBy,
	}

	if err := s.membershipRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("could not create membership: %w", err)
	}

	if s.transactionRepo != nil {
		_ = s.transactionRepo.Create(ctx, &transaction.Transaction{
			ReferenceCode: card.CardCode,
			Type:          "membership",
			Amount:        fee,
			ProcessedAt:   now,
			ProcessedBy:   issuedBy,
		})
	}

	return toMembershipResponse(card), nil
}

func (s *membershipService) GetMembershipByCode(ctx context.Context, code string) (*membership.MembershipResponse, error) {
	card, err := s.membershipRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return toMembershipResponse(card), nil
}

func (s *membershipService) GetAllMemberships(ctx context.Context) ([]membership.MembershipResponse, error) {
	cards, err := s.membershipRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return toMembershipResponses(cards), nil
}

func (s *membershipService) GetActiveMemberships(ctx context.Context) ([]membership.MembershipResponse, error) {
	cards, err := s.membershipRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	return toMembershipResponses(cards), nil
}

func toMembershipResponse(card *membership.MembershipCard) *membership.MembershipResponse {
	return &membership.MembershipResponse{
		ID:               card.ID,
		CardCode:         card.CardCode,
		PlateNumber:      card.PlateNumber,
		RegistrationDate: card.RegistrationDate,
		ExpirationDate:   card.ExpirationDate,
		IsActive:         card.IsActive(),
		BarcodeURL:       util.BarcodeURL(card.CardCode),
	}
}

func toMembershipResponses(cards []membership.MembershipCard) []membership.MembershipResponse {
	result := make([]membership.MembershipResponse, 0, len(cards))
	for i := range cards {
		result = append(result, *toMembershipResponse(&cards[i]))
	}
	return result
}
