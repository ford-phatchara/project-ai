package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/database"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/handlers"
)

func main() {
	// Initialize database connection
	database.InitDatabase()

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)

	// Start server on port 8080
	r.Run(":8080")
}
