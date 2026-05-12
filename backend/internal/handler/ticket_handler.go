package handler

import (
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/service"

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

// CreateTicket ออกตั๋ว (Check-in)
func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var req ticket.CreateTicketRequest
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

	res, err := h.service.CreateTicket(c.Context(), &req)
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

// CheckTicket สแกนออก (Check-out)
func (h *TicketHandler) CheckTicket(c *fiber.Ctx) error {
	var req ticket.CheckTicketRequest
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

	res, err := h.service.CheckTicket(c.Context(), &req)
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

// SearchTicket ค้นหาตั๋วจากทะเบียนรถ
func (h *TicketHandler) SearchTicket(c *fiber.Ctx) error {
	plateNumber := c.Query("plate_number")
	if plateNumber == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "plate_number is required",
		})
	}

	req := &ticket.SearchTicketRequest{PlateNumber: plateNumber}
	res, err := h.service.SearchTicket(c.Context(), req)
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

// GetIncome ดูรายได้ (daily / weekly / monthly)
func (h *TicketHandler) GetIncome(c *fiber.Ctx) error {
	period := c.Query("period", "daily")

	req := &ticket.IncomeRequest{Period: period}
	res, err := h.service.GetIncome(c.Context(), req)
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
