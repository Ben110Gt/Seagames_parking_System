package handler

import (
	"seagame/ticket/backend/internal/service"
	util "seagame/ticket/backend/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IncomeHandler struct{ service service.IncomeService }

func NewIncomeHandler(service service.IncomeService) *IncomeHandler { return &IncomeHandler{service} }

// GET /api/v1/admin/income/daily
func (h *IncomeHandler) Daily(c *fiber.Ctx) error {
	dateStr := c.Query("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.Fail("invalid date format, use YYYY-MM-DD"))
	}
	report, err := h.service.DailyReport(c.Context(), date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(report, ""))
}

// GET /api/v1/admin/income/weekly
func (h *IncomeHandler) Weekly(c *fiber.Ctx) error {
	startStr := c.Query("start", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.Fail("invalid date format, use YYYY-MM-DD"))
	}
	report, err := h.service.WeeklyReport(c.Context(), start)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(report, ""))
}

// GET /api/v1/admin/income/monthly
func (h *IncomeHandler) Monthly(c *fiber.Ctx) error {
	now := time.Now()
	year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(now.Year())))
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(now.Month()))))
	report, err := h.service.MonthlyReport(c.Context(), year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(util.Fail(err.Error()))
	}
	return c.JSON(util.OK(report, ""))
}
