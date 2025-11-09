package database

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todolistapi/models"
)

var Db *gorm.DB

func InitDB() {
	var err error

	// Use SQLite for local development, PostgreSQL for production
	if os.Getenv("DATABASE_URL") != "" {
		// Production (Render) - use PostgreSQL
		log.Println("Using PostgreSQL database...")
		Db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	} else {
		// Local development - use SQLite
		log.Println("Using SQLite database for local development...")
		Db, err = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto migrate tables
	err = Db.AutoMigrate(&models.User{}, &models.TodoItem{})
	if err != nil {
		log.Fatal("Error migrating table: ", err)
	}
	
	log.Println("Database connected and migrated successfully")
}
