package main

import (
	"fmt"
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
	// POST /register: Register a new user.
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	// POST /login: Log in and receive a JWT token.
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// protected routes
	protectedRouter := r.PathPrefix("/todos").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// POST /todos: Create a todo (authenticated).
	//protectedRouter.HandleFunc("/", handlers.CreateTodos).Methods("POST")
	protectedRouter.HandleFunc("/create", handlers.CreateTodos)
	// GET /todos: Get all todos belonging to the logged-in user.
	//protectedRouter.HandleFunc("/", handlers.GetTodos).Methods("GET")
	protectedRouter.HandleFunc("/list", handlers.GetTodos)
	// PUT /todos/:id: Update a specific todo (only if it belongs to the user).
	//protectedRouter.HandleFunc("/{id}", handlers.UpdateTodo).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", handlers.UpdateTodo)
	// DELETE /todos/:id: Delete a specific todo (only if it belongs to the user).
	//protectedRouter.HandleFunc("id/{id}", handlers.DeleteTodo).Methods("DELETE")
	protectedRouter.HandleFunc("/delete/{id}", handlers.DeleteTodo)

	// start the server
	fmt.Println("Server is running")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
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
