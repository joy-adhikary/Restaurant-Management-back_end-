package controllers

import (
	"context"
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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancle()

		result, err := tableCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure when all table information fatch"})
		}

		var alltables []bson.M

		if err := result.All(ctx, &alltables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure when all table information fatch 2"})
		}
		c.JSON(http.StatusOK, alltables)
	}
}

func GetTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var table models.Table

		tableId := c.Param("table_id")

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure wn listing item from table "})
			return
		}

		c.JSON(http.StatusOK, table)

	}
}

func UpdateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var table models.Table

		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not create table "})
			return
		}

		var updateObj primitive.D

		if table.Number_of_guests != nil {

			updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}
		if table.Table_number != nil {

			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_at})

		filter := bson.M{"table_id": tableId}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		) // tableItemCollection.UpdateOne(ctx,bson.M{"table_id": tableId},bson.D{{"$set", updateObj}, },&opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed at updating tableitems"})
		}
		// defer cancle()
		c.JSON(http.StatusOK, result)

	}
}

func CreateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not create table "})
			return
		}

		validationErr := validate.Struct(table)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, err := tableCollection.InsertOne(ctx, table)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table item not created "})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
