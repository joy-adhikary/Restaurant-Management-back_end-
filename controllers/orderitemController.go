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

type OrderItemPack struct {
	Table_id    *string
	order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

		var allOrderItems []bson.M

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur when we listing orderitem"})
		}

		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure wn listing item by orderitem"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.order_items { //order_items []models.OrderItem

			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_Item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_Price, 2)
			orderItem.Unit_Price = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)

		}

		// orderitem er majhe onk gula order thakbe .. mane akta slice tahkbe r akta table num tahkbe..silce er majhe onk gula item or [index] thkbe ..
		// ajonno amake sob  gula orderitem er kisu kisu data realtime update korty hobe
		// as like ID , created_at , updated_at to seijonno amake oi slice a loop calai akta akta kore value(item) niye update korty hobe
		//thn oi update kora slice index gula (akta full set of struct ) ke akta slice er majhe rakhtyci jeita orderitemtobeinserted

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancle()
		c.JSON(http.StatusOK, insertedOrderItems)

	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancle()
 
		var orderItem models.OrderItem

		   err:=c.BindJSON(&orderItem);err!=nil{
		     c.JSON(http.StatusInternalServerError,gin.H{"error":"error occure when update orderitem "})
		   }

		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}

		var Updateobj premitive.D

		if orderItem.Unit_Price != nil {
			Updateobj = append(Updateobj, bson.E{"unit_price": *&orderItem.Unit_Price})
		}
		if orderItem.Food_id != nil {
			Updateobj = append(Updateobj, bson.E{"food_id", *orderItem.Food_id})
		}
		if orderItem.Quantity != nil {
			Updateobj = append(Updateobj, bson.E{"quantity", *orderItem.Quantity})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		Updateobj = append(Updateobj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		) // orderItemCollection.UpdateOne(ctx,bson.M{"menu_id": menuId},bson.D{{"$set", updateObj}, },&opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed at updating orderitem"})
		}
		defer cancle()
		c.JSON(http.StatusOK, result)

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order by id"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)

	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}
