package routes

import (
	controller "Restaurant-Management-back_end/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/foods", controller.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controller.GetFood())
	incomingRoutes.POST("/foods", controller.CreateFood()) // create new item
	incomingRoutes.PATCH("/foods/:food_id", controller.UpdateFood())

}
