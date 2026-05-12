package handler

import (
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MembershipHandler struct {
	service  service.MembershipService
	validate *validator.Validate
}

func NewMembershipHandler(memberService service.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		service:  memberService,
		validate: validator.New(),
	}
}

// CreateMembership สร้าง/ต่ออายุ membership
func (h *MembershipHandler) CreateMembership(c *fiber.Ctx) error {
	var req ticket.CreateMembershipRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	res, err := h.service.CreateMembership(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

// CheckMembership ตรวจสอบสถานะ membership
func (h *MembershipHandler) CheckMembership(c *fiber.Ctx) error {
	plateNumber := c.Query("plate_number")
	if plateNumber == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "plate_number is required",
		})
	}

	req := &ticket.CheckMembershipRequest{PlateNumber: plateNumber}
	res, err := h.service.CheckMembership(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}
