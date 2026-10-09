# To-do List REST API

A REST API in Go for managing personal to-do lists, with user registration, JWT authentication and per-user access control on every to-do endpoint.

## What it does

- Users register and log in. Passwords are hashed with bcrypt before they are stored.
- Logging in returns a JWT that lasts 24 hours. Every `/todos` route sits behind authentication middleware that checks the token and reads the user ID from it.
- Each user can only see, update or delete their own to-dos. Every query is scoped to the logged-in user, so a request for someone else's to-do returns "not found".
- Each to-do has a title, description, completed flag, optional due date and priority (Medium by default).
- It uses SQLite for local development and PostgreSQL in production, through the GORM ORM, which also handles the table migrations.

## Endpoints

| Method | Route | Auth | Purpose |
|---|---|---|---|
| GET | `/` | No | Health check |
| POST | `/register` | No | Create a user |
| POST | `/login` | No | Log in and receive a JWT |
| POST | `/todos/create` | Yes | Create a to-do |
| GET | `/todos/list` | Yes | List your to-dos |
| PUT | `/todos/{id}` | Yes | Update one of your to-dos |
| DELETE | `/todos/delete/{id}` | Yes | Delete one of your to-dos |

Send the token as `Authorization: Bearer <token>` on protected routes.

## Tech

Go, Gorilla Mux, GORM, SQLite, PostgreSQL, golang-jwt, bcrypt

## Run it locally

```bash
export JWT_SECRET=choose-a-long-random-string
go run main.go
```

The server starts on port 8081 unless `PORT` is set. Without `DATABASE_URL` it creates a local SQLite database called `local.db`.

## Example

```bash
curl -X POST localhost:8081/register -d '{"first_name":"Fatima","last_name":"Sharif","email":"me@example.com","password":"secret"}'
curl -X POST localhost:8081/login -d '{"email":"me@example.com","password":"secret"}'
curl -X POST -H "Authorization: Bearer <token>" localhost:8081/todos/create -d '{"title":"Prep for interview","priority":"High"}'
curl -H "Authorization: Bearer <token>" localhost:8081/todos/list
```
