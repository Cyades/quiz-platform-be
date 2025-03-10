package routes

import (
	"quiz-platform/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tryout routes
		tryouts := v1.Group("/tryouts")
		{
			tryouts.GET("", controllers.GetAllTryouts)
			tryouts.GET("/:id", controllers.GetTryout)
			tryouts.POST("", controllers.CreateTryout)
			tryouts.PUT("/:id", controllers.UpdateTryout)
			tryouts.DELETE("/:id", controllers.DeleteTryout)

			// Helper route for options/filtering
			tryouts.GET("/options", controllers.GetTryoutOptions)
		}
	}

	return router
}
