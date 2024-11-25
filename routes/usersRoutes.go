package routes

import (
	"go_backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// User-related routes
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/register", controllers.Register)
		userGroup.PUT("/updateProfile", controllers.UpdateProfile)
		userGroup.PUT("/changePassword", controllers.ChangePassword)
		userGroup.GET("/getAll/:searchTerm", controllers.GetAll)
		userGroup.PUT("/toggleBlock/:userId", controllers.ToggleBlock)
		userGroup.GET("/getById/:userId", controllers.GetById)
		userGroup.PUT("/update", controllers.UpdateUser)
	}
}
