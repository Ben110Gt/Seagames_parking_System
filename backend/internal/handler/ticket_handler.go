package handler

import (
	"strings"

	"seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/service"
	util "seagame/ticket/backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	service  service.TicketService
	validate *validator.Validate
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	validate := validator.New()
	return &TicketHandler{
		service:  ticketService,
		validate: validate,
	}
}

// POST /api/v1/tickets
func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var req ticket.CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(util.Fail("invalid request body"))
	}
	if req.PlateNumber == "" {
		return c.Status(400).JSON(util.Fail("plate_number is required"))
	}
	req.IssuedBy = middleware.GetUserID(c)
	resp, err := h.service.CreateTicket(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.Status(201).JSON(util.OK(resp, "ticket created"))
}

// POST /api/v1/tickets/checkout
func (h *TicketHandler) Checkout(c *fiber.Ctx) error {
	var body struct {
		TicketCode string `json:"ticket_code"`
	}
	if err := c.BodyParser(&body); err != nil || body.TicketCode == "" {
		return c.Status(400).JSON(util.Fail("ticket_code is required"))
	}
	resp, err := h.service.CheckTicket(c.Context(), &ticket.CheckTicketRequest{
		TicketCode: body.TicketCode,
		CheckedBy:  middleware.GetUserID(c),
	})
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			return c.Status(404).JSON(util.Fail(msg))
		case strings.Contains(msg, "already been checked out"), strings.Contains(msg, "expired"):
			return c.Status(409).JSON(util.Fail(msg))
		default:
			return c.Status(400).JSON(util.Fail(msg))
		}
	}
	return c.Status(200).JSON(util.OK(resp, resp.Message))
}

// POST /api/v1/tickets/search
func (h *TicketHandler) SearchTicket(c *fiber.Ctx) error {
	var req ticket.SearchTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(util.Fail("invalid request body"))
	}
	resp, err := h.service.SearchTicket(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.Status(200).JSON(util.OK(resp, "ticket searched"))
}

// GET /api/v1/tickets/active
func (h *TicketHandler) GetActiveTickets(c *fiber.Ctx) error {
	resp, err := h.service.GetActiveTickets(c.Context())
	if err != nil {
		return c.Status(500).JSON(util.Fail(err.Error()))
	}
	return c.Status(200).JSON(util.OK(resp, "active tickets"))
}
