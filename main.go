package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terena-info/terena.godriver/gomgo"
	"github.com/terena-info/terena.godriver/middlewares"
	"github.com/terena-info/terena.godriver/response"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// type PaymenMethod struct {
// 	Title    string `validate:"required" form:"title" json:"title" bson:"title"`
// 	UserId   string `validate:"objectId" form:"user_id" json:"user_id" bson:"user_id"`
// 	Time     string `validate:"date" form:"time" json:"time" bson:"time"`
// 	DateTime string `validate:"datetime" form:"datetime" json:"datetime" bson:"datetime"`
// }

// func SanitizeRequest() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ctx.Next()
// 	}
// }

type User struct {
	gomgo.DefaultField `bson:",inline"`
	Email              string `validate:"required,email" json:"email" form:"email" bson:"email"`
	Password           string `json:"password" form:"password" bson:"password"`
	IsVerified         bool   `json:"is_verified,omitempty" form:"is_verified" bson:"is_verified,omitempty"`
	ProfileIcon        string `json:"profile_icon" form:"profile_icon" bson:"profile_icon"`
}

func main() {

	app := gin.Default()

	connector := gomgo.ConnectionOption{
		Host:     "mongodb+srv://bank:Bank211998Tsc_@cluster0.ih5kz.mongodb.net/?retryWrites=true&w=majority",
		Database: "terena_core",
		ReadRef:  readpref.Primary(),
		Timeout:  time.Second * 10,
		Context:  context.Background(),
	}
	connector.Connect().WithMessage(fmt.Sprintf("Database: %s", "terena_core"))

	app.Use(gin.CustomRecovery(middlewares.ErrorRecovery))

	// app.Use(SanitizeRequest())

	app.GET("/", func(ctx *gin.Context) {
		res := response.New(ctx)

		var user []interface{}
		hunter := []User{
			{
				Email:    "asdasdasd111",
				Password: "asdasdasdasdasd",
			},
			{
				Email:    "asdasdasd111",
				Password: "asdasdasdasdasd",
			},
		}

		for _, v := range hunter {
			user = append(user, v)
		}

		// query := gomgo.New(context.TODO(), "users")
		// user[0].ID = primitive.NewObjectID()

		// user[1].Email = "asdasdasd 1"
		// user[1].Password = "asdasdasd 1"
		// user[1].ID = primitive.NewObjectID()

		buldData(hunter)

		// var result User
		// query.InsertMany(user).Decode()

		// exist := query.FindOne(bson.M{"_id": utils.StringToObjectId("636235a775298f01db33fe12")}).Decode(&user).Exist()
		// if !exist {
		// 	gerrors.Panic(400, gerrors.E{Message: "NOT EXISTED DER"})
		// }

		res.Json(response.H{Data: user})
	})

	app.Run(":9009")

}

func buldData(a []interface{}) {
	b := make([]interface{}, len(a))
	for i := range a {
		b[i] = a[i]
	}

}
