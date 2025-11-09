package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gorilla/mux"
	"net/http"
	"todolistapi/database"
	"todolistapi/handlers"
	"todolistapi/middleware"
)

func main() {
	database.InitDB()

	// routes
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protectedRouter := r.PathPrefix("/todos").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	protectedRouter.HandleFunc("/create", handlers.CreateTodos).Methods("POST")
	protectedRouter.HandleFunc("/list", handlers.GetTodos).Methods("GET")
	protectedRouter.HandleFunc("/{id}", handlers.UpdateTodo).Methods("PUT")
	protectedRouter.HandleFunc("/delete/{id}", handlers.DeleteTodo).Methods("DELETE")

	// Add a simple health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Root endpoint hit!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "TodoList API is running!")
	}).Methods("GET")

	// Get port from environment (Render sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default for local development
	}

	// Start the server
	log.Printf("Server starting on port %s...", port)
	log.Println("Available endpoints:")
	log.Println("  GET  /")
	log.Println("  POST /register") 
	log.Println("  POST /login")
	log.Println("  POST /todos/create")
	log.Println("  GET  /todos/list")
	log.Println("  PUT  /todos/{id}")
	log.Println("  DELETE /todos/delete/{id}")
	
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
/*
Requirements:
Implement user registration and login functionality.✅
Use JWT authentication to protect the todo endpoints.
Create middleware that:
Authenticates requests using JWT.✅
STEPS

//1.create model ✅
//2.add model to authmigration in db.go ✅
//3.impliment handler function in api.go ✅
//4.declare the relationship in main.go ✅
//5.connect to database
Core Features:
POST /register: Register a new user. ✅
POST /login: Log in and receive a JWT token. ✅
POST /todos: Create a todo (authenticated).✅
GET /todos: Get all todos belonging to the logged-in user.✅
PUT /todos/:id: Update a specific todo (only if it belongs to the user).✅
DELETE /todos/:id: Delete a specific todo (only if it belongs to the user).✅

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI3MTMxMjgsInVzZXJfaWQiOjF9.Y_mhmbeauUtiZR05TR2t3dC2lnK4y7wbD7vH8-Avx6M"}



*/
