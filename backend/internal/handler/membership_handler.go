package handler

import (
	"seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/models/membership"
	"seagame/ticket/backend/internal/service"
	util "seagame/ticket/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type MembershipHandler struct {
	service service.MembershipService
}

func NewMembershipHandler(service service.MembershipService) *MembershipHandler {
	return &MembershipHandler{service}
}

func (h *MembershipHandler) CreateMembership(c *fiber.Ctx) error {
	var req membership.CreateMembershipRequest
	if err := c.BodyParser(&req); err != nil || req.PlateNumber == "" {
		return c.Status(400).JSON(util.Fail("plate_number is required"))
	}
	resp, err := h.service.CreateMembership(c.Context(), &req, middleware.GetUserID(c))
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.Status(201).JSON(util.OK(resp, "membership created"))
}

// GET /api/v1/memberships
func (h *MembershipHandler) GetAllMemberships(c *fiber.Ctx) error {
	resp, err := h.service.GetAllMemberships(c.Context())
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(resp, ""))
}

// GET /api/v1/memberships/active
func (h *MembershipHandler) GetActiveMemberships(c *fiber.Ctx) error {
	resp, err := h.service.GetActiveMemberships(c.Context())
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(resp, ""))
}

// GET /api/v1/memberships/:code
func (h *MembershipHandler) GetMembershipByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(400).JSON(util.Fail("code is required"))
	}

	resp, err := h.service.GetMembershipByCode(c.Context(), code)
	if err != nil {
		return c.Status(404).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(resp, ""))
}
