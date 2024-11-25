package routes

import (
	"go_backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupFoodsRouter(router *gin.Engine) {
	// Food-related routes
	foodGroup := router.Group("/api/foods")
	{
		foodGroup.GET("", controllers.GetAllFoods)
		foodGroup.GET("/search/:searchTerm", controllers.SearchFoods)
		foodGroup.GET("/tags", controllers.GetAllTags)
		foodGroup.GET("/tag/:tag", controllers.GetFoodsByTag)
		foodGroup.GET("/:foodId", controllers.GetFoodByID)
		foodGroup.DELETE("/:foodId", controllers.DeleteFood)
		foodGroup.PUT("/", controllers.UpdateFood)
		foodGroup.POST("/", controllers.AddFood)
	}
}
