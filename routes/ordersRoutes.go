package routes

import (
	"go_backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Order routes
	router.POST("/api/orders/create", controllers.CreateOrder)
	router.GET("/api/orders/newOrderForCurrentUser", controllers.GetNewOrderForCurrentUser)
	router.PUT("/api/orders/pay", controllers.Pay)
	router.GET("/api/orders/track/:orderId", controllers.TrackOrderById)
	router.GET("/api/orders/:state", controllers.GetAll)
	router.GET("/api/orders/allstatus", controllers.GetAllStatus)

	return router
}
