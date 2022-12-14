package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joy-adhikary/Restaurant-Management-back_end/database"
	"github.com/joy-adhikary/Restaurant-Management-back_end/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodcollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

//var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetFoods() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }

		// page, err := strconv.Atoi(c.Query("page"))
		// if err != nil || page < 1 { // page negative or 0 hole
		// 	page = 1 // 1st page dekhbo
		// }

		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// matchStage := bson.D{{"$match", bson.D{{}}}}                                                                     // match record by criteria in the db
		// groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}}}} // group all record by criteria .. for example name group .. sob  gula ke group korbe name er upr
		// projectStage:=                                                                                              // front end a ki dekhbo seita

		result, err := foodcollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "data fatch issue arise "})
			return

		}

		var allFoods []bson.M

		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allFoods)

	}
}

func GetFood() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		foodId := c.Param("food_id")

		var food models.Food

		err := foodcollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food) // foodCollection er food_id er songe match korbe (params er food_id) foodId
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in fatching data "})
		}
		c.JSON(http.StatusOK, food) //it will encode the row file into json format and show at web
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		foodId := c.Param("food_id")

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "food_name", Value: food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price}) // bson.E kontar majhe : konta rakhbo
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			if err != nil {
				msg := "menu was not found "
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu", Value: food.Price})
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.Updated_at})

		filter := bson.M{"food_id": foodId}
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodcollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error on update food item"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// foodId := c.Param("food_id")
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil { //catch the web json formar and decode it into normal form and push it into food
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationErr := validate.Struct(food) // validate if data correct or not

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": *food.Menu_id}).Decode(&menu) // bson.M{kontar songe : konta mathch korbo }

		if err != nil {
			msg := "menu not found"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		num := toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodcollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := "food not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	out := math.Pow(10, float64(precision))
	return float64(round(num * out))

}
