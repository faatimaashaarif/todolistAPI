package database

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todolistapi/models"
)

var Db *gorm.DB

func InitDB() {
	// Get database URL from environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	// Open database connection
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
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