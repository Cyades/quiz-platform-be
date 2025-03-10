package main

import (
	"fmt"
	"log"
	"os"
	"quiz-platform/config"
	"quiz-platform/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Connect to MongoDB
	config.ConnectDB()
	defer config.CloseDB()

	// Seed database with dummy data if empty
	config.SeedDummyData()

	// Set up router
	router := routes.SetupRouter()

	// Add a root route handler
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Quiz Platform API is running",
			"endpoints": []string{
				"/api/v1/tryouts",
				"/api/v1/tryouts/:id",
			},
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Server running on port %s...\n", port)
	fmt.Printf("Try accessing http://localhost:%s/api/v1/tryouts\n", port)
	router.Run(":" + port)
}
