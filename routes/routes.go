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
	
	// Root route for API health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Quiz Platform API is running",
			"endpoints": []string{
				"/api/v1/tryouts",
				"/api/v1/tryouts/:id",
				"/api/v1/tryouts/filter/options",
			},
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tryout routes
		tryouts := v1.Group("/tryouts")
		{
			tryouts.GET("", controllers.GetAllTryouts)
			tryouts.POST("", controllers.CreateTryout)
			
			// Helper route for options/filtering - must come before :id route to avoid conflict
			tryouts.GET("/filter/options", controllers.GetTryoutOptions)
			
			// Individual tryout routes with ID parameter
			tryouts.GET("/:id", controllers.GetTryout)
			tryouts.PUT("/:id", controllers.UpdateTryout)
			tryouts.DELETE("/:id", controllers.DeleteTryout)
		}
	}

	return router
}