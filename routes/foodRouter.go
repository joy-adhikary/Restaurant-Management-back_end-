package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/controllers"
)

func FoodRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood()) // create new item
	incomingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())

}
