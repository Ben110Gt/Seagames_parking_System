package database

import (
	"fmt"
	"log"
	"os"
	"seagame/ticket/backend/internal/models/membership"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/models/transaction"
	usermodel "seagame/ticket/backend/internal/models/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// ลอง load .env จากหลาย path
	envPaths := []string{".env", "../.env", "../../.env"}
	loaded := false
	for _, p := range envPaths {
		if err := godotenv.Load(p); err == nil {
			loaded = true
			log.Printf("Loaded .env from: %s", p)
			break
		}
	}
	if !loaded {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Database connection
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("Database connected successfully")

	//Migrate the schema
	if err := db.AutoMigrate(
		&usermodel.User{},
		&ticket.Ticket{},
		&membership.MembershipCard{},
		&transaction.Transaction{},
	); err != nil {
		log.Printf("Warning: AutoMigrate error: %v", err)
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
