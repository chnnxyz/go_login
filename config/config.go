package config

import (
	"cyberia_auth/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	// Load from env or hardcode if testing
	host := os.Getenv("DB_HOST")         // e.g. "mydb.abcdefg12345.us-east-1.rds.amazonaws.com"
	user := os.Getenv("DB_USER")         // e.g. "postgres"
	password := os.Getenv("DB_PASSWORD") // e.g. "mypassword"
	dbname := os.Getenv("DB_NAME")       // e.g. "tickets"
	port := os.Getenv("DB_PORT")         // e.g. "5432"
	sslmode := os.Getenv("DB_SSLMODE")   // usually "require" or "disable"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
	)
	log.Println(dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Role{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	log.Println("Connected to PostgreSQL auth database.")
}
