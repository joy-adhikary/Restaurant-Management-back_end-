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
			defer cancle()
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
				defer cancle()
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

		err := c.BindJSON(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure when update orderitem"})
		}

		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}

		var Updateobj primitive.D

		if orderItem.Unit_Price != nil {
			Updateobj = append(Updateobj, bson.E{Key: "unit_price", Value: orderItem.Unit_Price})
		}
		if orderItem.Food_id != nil {
			Updateobj = append(Updateobj, bson.E{Key: "food_id", Value: *orderItem.Food_id})
		}
		if orderItem.Quantity != nil {
			Updateobj = append(Updateobj, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		Updateobj = append(Updateobj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: Updateobj},
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

	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancle()

	//mongo work with json data

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}} // order item er (order id=id eita pass hocche ) er songe match korbo orderitem.order_id
	// matchstage a akta key diye oi key er sob record fatch kora jai
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}}}}
	//from => koi theke dekhbo (food collection er majhe ), amr local field konta , amr foreign field konta , ki hisabe dkehbo ei datagula
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	// normal array ke handel korty pare nah mongo cant perform any oparetion on it .. thats why we need to unwind it
	// agerline er food ke ei jaigai path a set korbo

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {Key: "as", Value: "order"}}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "order.table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}}}}
	// eijaigai order.table_id dicche karon join korar por amr akn orderitem r order er sob coloum/attribute/ values amr order er majhe store hocche
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "total_count", Value: 1},
			{Key: "food_name", Value: "$food.name"},
			{Key: "food_image", Value: "$food.food_image"},
			{Key: "table_number", Value: "$table.table_number"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "price", Value: "$food.price"},
			{Key: "quantity", Value: 1},
		}}}
	// manage the front end// what goes to the front end

	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "order_id", Value: "$order_id"}, {Key: "table_id", Value: "$table_id"}, {Key: "table_number", Value: "$table_number"}}}, {Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "order_items", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	//data gula ke group korbe akta perameter er upr base kore

	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancle()

	return OrderItems, err

}
