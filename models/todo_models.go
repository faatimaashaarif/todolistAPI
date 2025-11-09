package models

import (
	"gorm.io/gorm"
	"time"
)
//USE GORM - CREATE RECORD
type User struct {
	gorm.Model

	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Tasks      []TodoItem `gorm:"foreignKey:UserID"`
}

type TodoItem struct {
	gorm.Model

	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority" gorm:"default:'Medium'"`

	//Existing User foreign key
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"-"`

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
