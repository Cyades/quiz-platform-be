package main

import (
	"fmt"
	"log"
	"os"
	"quiz-platform/config"
	"quiz-platform/routes"
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

	// Set up router
	router := routes.SetupRouter()

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
