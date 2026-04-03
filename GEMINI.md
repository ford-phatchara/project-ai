## Project Overview
This is a Golang backend API project.
- Language: Go (latest stable version)
- Architecture: Clean architecture (handler, service, repository)
- Framework: Gin
- ORM: GORM
- Database: SQLite (for local development)

## Project Structure
- cmd/server: entry point
- internal/handlers: HTTP handlers
- internal/services: business logic
- internal/repositories: database access
- internal/models: data models
- internal/database: DB connection setup

## Coding Guidelines
- Keep code simple and readable
- Follow Go conventions (gofmt, idiomatic Go)
- Use dependency injection where possible
- Avoid global variables
- Write small, focused functions

## Coding Guidelines
- Keep code simple and readable
- Follow Go conventions (gofmt, idiomatic Go)
- Use dependency injection where possible
- Avoid global variables
- Write small, focused functions

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
- Install dependencies: go mod tidy
- Run server: go run cmd/server/main.go
- Default port: 8080

## Testing
- Write unit tests for services
- Use Go testing package

## Constraints
- Do not introduce unnecessary libraries
- Do not over-engineer
- Keep it beginner-friendly

## Instructions for AI
- When adding new features, follow the existing project structure
- Always separate handler, service, and repository layers
- Provide clear explanations when making changes
- Prefer simple and maintainable solutions over complex ones