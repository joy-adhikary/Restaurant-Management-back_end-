package main

import (
	"os"

	"github.com/gin-gonic/gin"
	middleware "github.com/joy-adhikary/Restaurant-Management-back_end/middleware"
	routes "github.com/joy-adhikary/Restaurant-Management-back_end/routes"
)

//var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New() //new gin router created
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication()) // check user is authenticat or not if yes then can use the router

	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.TableRoutes(router)

	router.Run(":" + port) // port a run korbe serveandlisten er mto

}

// {
// 	"First_name":"joy",
// 	"Last_name":"adhikary",
// 	"Password":"112233",
// 	"Email":"joyadhikary071@gmail.com",
// 	"Phone":"01517144672",
// 	"User_type":"ADMIN"
//   }
