package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/database"
	"github.com/joy-adhikary/Restaurant-Management-back_end/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		result, err := menuCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure when data fatch"})
		}

		var allMenus []bson.M

		if err = result.All(ctx, &allMenus); err != nil { //result er sob data allMenus er majhe json bind kore dibe
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)

	}
}

func GetMenu() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		menuId := c.Param("Menu_id")
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"Menu_id": menuId}).Decode(&menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurs in fatching data"})
		}
		c.JSON(http.StatusOK, menu)

	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var menuId = c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_date != nil && menu.End_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ops error"})
				defer cancle()
				return
			}
		}

		updateObj = append(updateObj, bson.E{"start_date", menu.Start_date})
		updateObj = append(updateObj, bson.E{"end_date", menu.End_date})

		if menu.Name != " " {
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}
		if menu.Catagory != " " {
			updateObj = append(updateObj, bson.E{"catagory", menu.Catagory})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		//menu.Created_at,_=time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		) // menuCollection.UpdateOne(ctx,bson.M{"menu_id": menuId},bson.D{{"$set", updateObj}, },&opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed "})
		}
		defer cancle()
		c.JSON(http.StatusOK, result)
	}
}

func CreateMenu() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validatorError := validate.Struct(menu)

		if validatorError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatorError.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, inserterr := menuCollection.InsertOne(ctx, menu)

		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error when inserting "})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}
