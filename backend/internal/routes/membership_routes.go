package routes

import (
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/handler"
	middleware "seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/repository"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupMembershipRoutes(app *fiber.App) {
	db := database.GetDB()
	memberRepo := repository.NewMembershipRepository(db)
	memberService := service.NewMembershipService(memberRepo)
	memberHandler := handler.NewMembershipHandler(memberService)

	// ตรวจสอบ membership (public)
	app.Get("/membership/check", memberHandler.CheckMembership)

	// สร้าง/ต่ออายุ membership (ต้อง auth)
	member := app.Group("/membership")
	member.Use(middleware.JWTMiddleware())
	member.Post("/", memberHandler.CreateMembership)
}
