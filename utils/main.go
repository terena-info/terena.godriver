package utils

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsIn(value string, list []string) bool {
	for _, val := range list {
		if value == val {
			return true
		}
	}
	return false
}

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return err.(validator.ValidationErrors)[0]
	}
	return nil
}

func StringToObjectId(value string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		panic("500::Failed to convert string to object id")
	}
	return objectId
}
