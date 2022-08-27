package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/database"
	"github.com/joy-adhikary/Restaurant-Management-back_end/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	}
}

func CreateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var table models.Table

		tableId := c.Param("table_id")

	}
}
