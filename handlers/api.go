package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"todolistapi/database"
	"todolistapi/middleware"
	"todolistapi/models"
	"todolistapi/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// collect the details of the user as request body
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the user already exist
	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Hash the password
	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}

	req.Password = HashPassword

	// add the user to the database
	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError) // Changed status to 500
		return
	}

	// send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User created successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	// decode the request the request body
	var login models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the user exists
	var user models.User
	err = database.Db.Where("email = ?", login.Email).First(&user).Error
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized) // Consolidated errors
		return
	}

	// check if password matches what we have in our database
	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// generating a token
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	// send a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func CreateTodos(w http.ResponseWriter, r *http.Request) { //(POST/todos/create)
	var createTodos models.TodoItem
	err := json.NewDecoder(r.Body).Decode(&createTodos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	createTodos.UserID = userID

	if err := database.Db.Create(&createTodos).Error; err != nil {
		http.Error(w, "Unable to create todo item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createTodos)
}

func GetTodos(w http.ResponseWriter, r *http.Request) { //(GET/todos/list)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var todos []models.TodoItem
	err := database.Db.Where("user_id = ?", userID).Find(&todos).Error
	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	// Send the list of todos back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) { //(PUT/todos/{id})
	var updateTodo models.TodoItem
	err := json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	todoIDStr := utils.GetURLParam(r, "id")
	if todoIDStr == "" {
		http.Error(w, "Missing todo ID in URL", http.StatusBadRequest)
		return
	}

	var existingTodo models.TodoItem

	// from db find the existing Todo item and ensure it belongs to the user
	db := database.Db.Where("user_id = ?", userID).First(&existingTodo, todoIDStr)

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Todo item not found or unauthorized", http.StatusNotFound)
		} else {
			fmt.Printf("Database error finding todo: %v\n", db.Error)
			http.Error(w, "Failed to retrieve todo item", http.StatusInternalServerError)
		}
		return
	}

	// check correct ID/UserID are preserved
	updateTodo.ID = existingTodo.ID
	updateTodo.UserID = userID // Ensure UserID is preserved and correct

	// save the updated item
	if err := database.Db.Save(&updateTodo).Error; err != nil {
		fmt.Printf("Database error updating todo: %v\n", err)
		http.Error(w, "Failed to update todo item", http.StatusInternalServerError)
		return
	}

	// send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateTodo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) { //(DELETE/todos/delete/{id})
	todoIDStr := utils.GetURLParam(r, "id")
	if todoIDStr == "" {
		http.Error(w, "Missing todo ID in URL", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	db := database.Db.Where("user_id = ?", userID).Delete(&models.TodoItem{}, todoIDStr)

	if db.Error != nil {
		fmt.Printf("Database error deleting todo: %v\n", db.Error)
		http.Error(w, "Failed to delete todo item", http.StatusInternalServerError)
		return
	}

	if db.RowsAffected == 0 {
		http.Error(w, "Todo item not found or you do not have permission", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
