package database

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todolistapi/models"
)

var Db *gorm.DB

func InitDB() {
	var connStr string
	
	// Try DATABASE_URL first (for Render)
	connStr = os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fall back to individual variables (for local development)
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if user == "" {
			user = "postgres"
		}
		if dbname == "" {
			dbname = "todolistapi"
		}
		
		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)
	}

	if connStr == "" {
		log.Fatal("Database connection string is not set")
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