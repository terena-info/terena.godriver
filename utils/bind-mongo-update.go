package utils

import (
	"encoding/json"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BindUpdate(data interface{}) primitive.D {
	v := reflect.ValueOf(data)
	typeOfS := v.Type()

	result := bson.D{{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())}}

	for i := 0; i < v.NumField(); i++ {
		f := typeOfS.Field(i)
		if f.Name != "DefaultField" {
			field := typeOfS.Field(i).Tag.Get("bson")
			val := v.Field(i).Interface()
			result = append(result, bson.E{Key: field, Value: val})
		}
	}
	return result
}

func BindCreate(m any) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(m)
	json.Unmarshal(inrec, &inInterface)
	inInterface["_id"] = primitive.NewObjectID()
	inInterface["created_at"] = primitive.NewDateTimeFromTime(time.Now())
	inInterface["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	inInterface["is_active"] = true
	re, _ := json.Marshal(inInterface)
	json.Unmarshal(re, &m)
}
