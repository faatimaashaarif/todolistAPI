package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"todolistapi/models"
)

var Db *gorm.DB

func InitDB() {
	// open the database
	connStr := "user=postgres host=localhost password=password dbname=todolistapi port=5432 sslmode=disable"

	var err error
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	err = Db.AutoMigrate(&models.User{}, &models.TodoItem{})
	if err != nil {
		log.Fatal("error migrating table", err)
	}
}
