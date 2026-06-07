package routes

import (
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/handler"
	middleware "seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/repository"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetIncomeRoutes(app *fiber.App) {
	db := database.GetDB()
	transactionRepo := repository.NewTransactionRepository(db)
	incomeService := service.NewIncomeService(transactionRepo)
	incomeHandler := handler.NewIncomeHandler(incomeService)

	admin := app.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.RoleMiddleware("Owner"))
	admin.Get("/income/daily", incomeHandler.Daily)
	admin.Get("/income/weekly", incomeHandler.Weekly)
	admin.Get("/income/monthly", incomeHandler.Monthly)
}
