package main

import (
	"github.com/gin-gonic/gin"
	"github.com/terena-info/terena.godriver/binding"
	"github.com/terena-info/terena.godriver/middlewares"
	"github.com/terena-info/terena.godriver/response"
)

type PaymenMethod struct {
	Title    string `validate:"required" form:"title" json:"title" bson:"title"`
	UserId   string `validate:"objectId" form:"user_id" json:"user_id" bson:"user_id"`
	Time     string `validate:"date" form:"time" json:"time" bson:"time"`
	DateTime string `validate:"datetime" form:"datetime" json:"datetime" bson:"datetime"`
}

func SanitizeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func main() {
	app := gin.Default()

	app.Use(gin.CustomRecovery(middlewares.ErrorRecovery))

	app.Use(SanitizeRequest())

	app.POST("/", func(ctx *gin.Context) {
		res := response.New(ctx)

		var admin PaymenMethod
		ctx.ShouldBind(&admin)

		validate := binding.New(admin)
		validate.ValidateStruct().RunError(&binding.RunErrorOption{})

		res.Json(response.H{})
	})

	app.Run(":9009")
}
