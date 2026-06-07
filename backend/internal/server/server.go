package server

import (
	"log"
	"os"
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func NewServer() {

	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file Server")
	}

	app := fiber.New(fiber.Config{
		AppName: "Seagames Parking System v1.0",
	})

	// Connect to the database
	database.ConnectDatabase()

	app.Use(logger.New())

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Seagames Parking System API",
		})
	})

	// Register all routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.TicketRoutes(app)
	routes.SetupMembershipRoutes(app)
	routes.SetIncomeRoutes(app)

	// Start the server
	port := os.Getenv("APP_PORT")

	app.Listen("0.0.0.0:" + port)
}
