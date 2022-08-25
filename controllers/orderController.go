package controllers

import (
	"context"
	"fmt"
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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var allorders []bson.M

		result, err := orderCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "did not find any orders"})
			return

		}
		if err = result.All(ctx, &allorders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "did not find any orders"})
			return
		}
		c.JSON(http.StatusOK, allorders)

	}
}

func GetOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		orderId := c.Param("order_id")

		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "there is no order id whcih is requested"})
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur on updation "})
			return
		}

		orderId := c.Param("order_id")

		var updateObj primitive.D

		if order.Table_id != nil {
			err := orderCollection.FindOne(ctx, bson.M{"table_id": orderId}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "table id not founded in the database"})
				return

			}
			updateObj = append(updateObj, bson.E{"table_id", table.Table_id})
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		filter := bson.M{"order_id": orderId}

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error when the update data is fatch "})
			return
		}
		defer cancle()
		c.JSON(http.StatusOK, result)
	}
}

func CreateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var table models.Table
		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error at create oparetion "})
			return
		}

		validateErr := validate.Struct(order)

		if validateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr.Error()})
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"table_id": *order.Table_id}).Decode(table)

		if err != nil {
			msg := fmt.Sprintf("tabel not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			msg := fmt.Sprintf("order not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func OrderItemOrderCreator(order models.Order) string {

	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancle()
	return order.Order_id

}
