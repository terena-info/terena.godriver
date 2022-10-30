package main

import (
	"context"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/terena-info/terena.godriver/gomgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Name string `json:"name" bson:"name"`
}

type PaymenMethod struct {
	Title string `json:"title" bson:"title"`
	Icon  string `json:"icon" bson:"icon"`
	Code  int    `json:"code" bson:"code"`
}

type OrderStatus struct {
	NameLA string `json:"name_la" bson:"name_la"`
	NameEN string `json:"name_en" bson:"name_en"`
}

type Order struct {
	UserID         string         `json:"user_id" bson:"user_id"`
	OrderNumber    int            `json:"order_number" bson:"order_number"`
	PaymenMethodID []PaymenMethod `json:"payment_method_id" bson:"payment_method_id"`
	OrderStatusID  []OrderStatus  `json:"order_status_id" bson:"order_status_id"`
}

func main() {
	db := gomgo.ConnectionOption{
		Host:     "mongodb+srv://paymart:LprAaM49lzuKMU6y@paymart.wwlmk.mongodb.net/development?retryWrites=true&w=majority",
		Database: "development",
		Timeout:  time.Second * 10, // Wait for 10s to success connect
		ReadRef:  readpref.Primary(),
		Context:  context.Background(),
	}
	db.Connect().WithMessage("Database is connected")

	// userId, _ := primitive.ObjectIDFromHex("6247c4389bf13408f7cfcc99")

	// var user []User
	// query.FindById(userId).Decode(&user).ErrorMessage("400::validate_user_id")
	app := gin.Default()

	app.GET("/", func(ctx *gin.Context) {
		query := gomgo.New(context.TODO(), "order")

		var orders []bson.M
		query.Lookup(&gomgo.LookupOption{
			From:         "payment_method",
			LocalField:   "payment_method_id",
			ForeignField: "_id",
			As:           "payment_method_id",
		}).Lookup(&gomgo.LookupOption{
			From:         "order_status",
			LocalField:   "order_status_id",
			ForeignField: "_id",
			As:           "order_status_id",
		}).AddPipeline(bson.M{
			"$lookup": bson.M{
				"from":         "seller",
				"localField":   "seller_id",
				"foreignField": "_id",
				"as":           "seller_id",
			},
		}).Unwind(&gomgo.UnwindOption{Path: "seller_id"}).AutoBindQuery(&gomgo.BindConfig{
			Context:       ctx,
			SearchFields:  []string{"order_number"},
			ForcePaginate: true,
		}).Select("order_number", "is_offline_order").Decode(&orders)

		ctx.JSON(200, orders[0])
	})

	app.Run(":9009")
}
