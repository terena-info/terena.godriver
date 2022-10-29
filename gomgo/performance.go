package gomgo

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOneIndex(key string, value int) func() {
	// Define index interface
	indexModel := mongo.IndexModel{
		Keys: bson.E{
			Key:   key,
			Value: value,
		},
	}

	// Create Index
	MongoInstance.Collection("User").Indexes().CreateOne(context.TODO(), indexModel)

	// return remove index after used or ignore
	return func() {
		MongoInstance.Collection("User").Indexes().DropOne(context.TODO(), fmt.Sprintf("%s_%s", key, strconv.Itoa(value)))
	}
}
