package routes

import (
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/handler"
	middleware "seagame/ticket/backend/internal/middleware"
	"seagame/ticket/backend/internal/repository"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository(database.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	admin := app.Group("/admin/users")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.RoleMiddleware("Owner"))
	admin.Get("/:id", userHandler.GetUserByID)
	admin.Get("/", userHandler.GetAllUsers)
	admin.Put("/:id", userHandler.UpdateUser)
	admin.Delete("/:id", userHandler.DeleteUser)
}
