package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/database"
	// "github.com/phatcharasangsuphap/gemlni-cli-backend/internal/handlers"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	// Initialize database connection
	database.InitDatabase()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define routes
	// r.GET("/users", handlers.GetUsers)
	// r.GET("/users/:id", handlers.GetUserById)
	// r.POST("/users", handlers.CreateUser)
	// r.PUT("/users/:id", handlers.UpdateUser)

	// Start server on port 8080
	r.Run(":8080")
}
