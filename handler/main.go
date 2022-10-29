package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/terena-info/terena.godriver/response"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerFunc struct {
	Error      error
	Message    string
	ErrorCode  string
	StatusCode int
	Data       interface{}
}

type HandlerContext struct {
	Context        *gin.Context
	SessionContext mongo.SessionContext
}

// Handler
func SessionHandler(handler func(HandlerContext), DBClient *mongo.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := response.New(ctx)
		handlerOption := HandlerContext{Context: ctx}
		session, err := DBClient.StartSession()
		if err != nil {
			res.Panic(http.StatusInternalServerError, response.HErr{Message: err.Error(), ErrorCode: "5000"})
		}
		defer session.EndSession(context.Background())
		session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
			handlerOption.SessionContext = sessCtx
			handler(handlerOption)
			return nil, nil
		})
	}
}

func Error(statusCode int, errorMessage string, errorCode string) HandlerFunc {
	return HandlerFunc{Error: errors.New(errorMessage), ErrorCode: errorCode, StatusCode: statusCode}
}
