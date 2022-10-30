package binding

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	gerrors "github.com/terena-info/terena.godriver/gerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RunErrorOption struct {
	StatusCode int
}

type _SchemaInterface interface {
	ValidateStruct() _Schema
	FirstError() string
	IsEmpty() int
	RunError(*RunErrorOption)
}

type _Schema struct {
	Schema interface{}
	errors []string
}

func objectId(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return primitive.IsValidObjectID(value)
}

func dateOnly(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

func dateTime(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse("2006-01-02 15:04:05", value)
	return err == nil
}

func (schema _Schema) ValidateStruct() _Schema {
	validate := validator.New()
	validate.RegisterValidation("objectId", objectId)
	validate.RegisterValidation("date", dateOnly)
	validate.RegisterValidation("datetime", dateTime)

	if err := validate.Struct(schema.Schema); err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			// Set tag value
			var errors string

			tagSet, _ := reflect.TypeOf(schema.Schema).FieldByName(v.StructField())
			errors = tagSet.Tag.Get("msg")

			if errors == "" {
				errors = fmt.Sprintf("validate_%s", tagSet.Tag.Get("json"))
			}

			if errors == "validate_" {
				errors = fmt.Sprintf("validate_%s", tagSet.Tag.Get("form"))
			}

			if errors == "validate_" {
				errors = fmt.Sprintf("validate_%s", strings.ToLower(v.StructField()))
			}

			schema.errors = append(schema.errors, errors)
		}
	}

	return schema
}

func (schema _Schema) IsEmpty() int {
	return len(schema.errors)
}

func (schema _Schema) FirstError() string {
	if schema.IsEmpty() > 0 {
		return schema.errors[0]
	}
	return ""
}

func (schema _Schema) RunError(opts *RunErrorOption) {
	if opts.StatusCode == 0 {
		opts.StatusCode = http.StatusBadRequest
	}
	if schema.FirstError() != "" {
		gerrors.Panic(opts.StatusCode, gerrors.E{Message: schema.FirstError(), ErrorCode: "4000"})
	}
}

func New(input interface{}) _SchemaInterface {
	var v _SchemaInterface = _Schema{Schema: input}
	return v
}
