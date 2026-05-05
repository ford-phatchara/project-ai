package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/database"
	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/routes"
)

// @title Durian Farm Management API
// @version 1.0
// @description API for managing durian farm plots, buyers, sales, maintenance logs, and dashboard analytics.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-User-ID
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register API routes
	routes.RegisterRoutes(r, database.DB)

	// Start server on port 8080
	r.Run(":8080")
}
