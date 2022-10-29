package main

import (
	"context"
	"fmt"
	"godriver/gomgo"
	"time"

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

	query := gomgo.New(context.TODO(), "order")

	// var user []User
	// query.FindById(userId).Decode(&user).ErrorMessage("400::validate_user_id")

	var orders []map[string]interface{}

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
	}).Select(
		"order_number",
		"payment_method_id",
	).Decode(&orders)

	fmt.Println(orders[0]["order_number"])
}
