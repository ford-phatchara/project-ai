# Fullstack User Management Project

This project features a Golang REST API backend and an Angular frontend.

## Project Structure

- **backend/**: Golang REST API using Gin and GORM.
  - `cmd/server/`: Main application entry point.
  - `internal/database/`: Database initialization and connection (SQLite).
  - `internal/handlers/`: API request handlers.
  - `internal/models/`: Database models.
- **frontend/**: Angular application using Standalone Components and RxJS.
  - `src/app/core/`: Services for API communication.
  - `src/app/features/`: UI components (User list, User form).
  - `src/app/models/`: Frontend data models.

---

## How to Run

### 1. Prerequisites
- [Go](https://go.dev/doc/install) (1.20+)
- [Node.js](https://nodejs.org/) (latest LTS)
- [npm](https://www.npmjs.com/) (latest)
- [Angular CLI](https://angular.dev/tools/cli) (optional, `npx` will be used if not installed)

### 2. Run the Backend
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```
4. The backend will be available at `http://localhost:8080`.

### 3. Run the Frontend
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```
4. Open your browser and go to `http://localhost:4200`.

---

## API Endpoints (Backend)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/users` | Fetch all users |
| GET    | `/users/:id` | Fetch a single user by ID |
| POST   | `/users` | Create a new user |
| PUT    | `/users/:id` | Update an existing user |

---

## Features
- **Frontend**: Angular standalone components, reactive forms, RxJS (Observable/pipe/map), and environment-based configuration.
- **Backend**: Gin framework, GORM (SQLite), and CORS middleware configured for `http://localhost:4200`.
- **Database**: SQLite database (`backend/gorm.db`) with auto-migration of the User model.
