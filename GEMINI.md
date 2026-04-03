## Project Overview
This is a fullstack project with:
- Backend: Golang REST API
- Frontend: Angular (latest)

## Project Structure
- backend:
  - cmd/server: entry point
  - internal/handlers: HTTP handlers
  - internal/services: business logic
  - internal/repositories: database access
  - internal/models: data models
  - internal/database: DB connection setup

- frontend:
  - src/app: Angular application
  - src/app/services: API services
  - src/app/components: UI components
  - src/app/models: frontend models

## Backend Stack
- Language: Go (latest stable version)
- Framework: Gin
- ORM: GORM
- Database: SQLite (local), can be upgraded later

## Frontend Stack
- Angular (latest)
- TypeScript
- RxJS
- Angular CLI

## Communication
- Frontend calls backend via REST API
- Default backend URL: http://localhost:8080
- Use JSON for request/response

## Coding Guidelines (Backend)
- Keep code simple and readable
- Follow Go conventions (gofmt, idiomatic Go)
- Use dependency injection where possible
- Avoid global variables
- Write small, focused functions

## Coding Guidelines (Frontend)
- Use Angular best practices
- Use services for API calls
- Use reactive programming (RxJS)
- Keep components small and reusable
- Separate logic from UI

## API Guidelines
- Use RESTful conventions
- JSON request/response only
- Proper HTTP status codes (200, 201, 400, 404, 500)
- Validate request input

## Database
- Use GORM for ORM
- Auto migrate models
- Keep models simple
- Use SQLite for local development

## Run Instructions

### Backend
- Install dependencies: `go mod tidy`
- Run server: `go run cmd/server/main.go`
- Default port: 8080

### Frontend
- Install dependencies: `npm install`
- Run app: `ng serve`
- Default port: 4200

## Testing
- Backend: use Go testing package
- Frontend: use Angular testing tools (Karma/Jasmine)

## Constraints
- Do not introduce unnecessary libraries
- Do not over-engineer
- Keep it beginner-friendly

## Instructions for AI
- When adding new features:
  1. Create API endpoint in backend
  2. Create service in Angular to call API
  3. Connect UI with API
- Always separate handler, service, repository layers in backend
- Keep frontend and backend loosely coupled
- Provide clear explanations when making changes
- Prefer simple and maintainable solutions over complex ones

## Example Tasks for AI
- Generate CRUD API in Go
- Generate Angular service to call API
- Generate Angular component for displaying data
- Connect frontend form to backend API