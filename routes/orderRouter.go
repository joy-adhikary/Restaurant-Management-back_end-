package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())

}
