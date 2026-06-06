package routes

import (
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/handler"
	middleware "seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/repository"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func TicketRoutes(app *fiber.App) {
	db := database.GetDB()

	tickRepo := repository.NewTicketRepository(db)
	memberRepo := repository.NewMembershipRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	tickService := service.NewTicketService(tickRepo, memberRepo, transactionRepo)
	ticketHandler := handler.NewTicketHandler(tickService)

	// app.Post("/tickets/create", ticketHandler.CreateTicket)

	// app.Post("/tickets/check", ticketHandler.CheckTicket)

	// app.Get("/tickets/search", ticketHandler.SearchTicket)

	admin := app.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.RoleMiddleware("Owner"))
	admin.Post("/create-ticket", ticketHandler.CreateTicket)
	admin.Post("/check-ticket", ticketHandler.Checkout)
	// admin.Get("/tickets/search", ticketHandler.SearchTicket)

	// Dashboard รายได้
	dashboard := app.Group("/dashboard")
	dashboard.Use(middleware.JWTMiddleware())
	dashboard.Use(middleware.RoleMiddleware("Owner"))
	// dashboard.Get("/income", ticketHandler.GetIncome)
}
