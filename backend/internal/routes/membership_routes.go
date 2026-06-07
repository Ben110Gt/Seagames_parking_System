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

	memRepo := repository.NewMembershipRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	memService := service.NewMembershipService(memRepo, transactionRepo)
	memHandler := handler.NewMembershipHandler(memService)

	admin := app.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.RoleMiddleware("Owner"))
	admin.Post("/create-membership", memHandler.CreateMembership)     // ✅
	admin.Get("/memberships", memHandler.GetAllMemberships)           // ✅
	admin.Get("/memberships/:code", memHandler.GetMembershipByCode)   // ✅
	admin.Get("/memberships/active", memHandler.GetActiveMemberships) // ❌

}
