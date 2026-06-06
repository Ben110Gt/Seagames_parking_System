package database

import (
	"log"

	"seagame/ticket/backend/internal/models/membership"
	"seagame/ticket/backend/internal/models/ticket"
	"seagame/ticket/backend/internal/models/user"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&user.User{},
		&ticket.Ticket{},
		&membership.MembershipCard{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database migration completed")
}
