# Golang Backend Project

This is a simple Golang backend project using Gin and GORM with SQLite.

## Project Structure

- `cmd/server/`: Main application entry point.
- `internal/database/`: Database initialization and connection.
- `internal/handlers/`: API request handlers.
- `internal/models/`: Database models.

## How to Run

1. Make sure you have Go installed (Go 1.26+ recommended).
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```
4. The server will start on `http://localhost:8080`.

## API Endpoints

### 1. Get all users
**GET** `/users`
```bash
curl http://localhost:8080/users
```

### 2. Create a user
**POST** `/users`
```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john@example.com"}'
```

### 3. Update a user
**PUT** `/users/:id`
```bash
curl -X PUT http://localhost:8080/users/1 \
     -H "Content-Type: application/json" \
     -d '{"name": "John Updated"}'
```

## Features

- **SQLite Database**: A simple file-based database (`gorm.db`) will be created automatically.
- **Auto Migration**: Database schema is automatically updated on server start.
- **Gin Framework**: Lightweight and fast HTTP web framework.
- **GORM**: Powerful Object-Relational Mapper for Go.
