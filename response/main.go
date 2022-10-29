package response

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type HErr struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type ResponseInterface interface {
	Json(H)
	Panic(int, HErr)
}

type Response struct {
	Ctx *gin.Context
}

func New(ctx *gin.Context) ResponseInterface {
	var response ResponseInterface = Response{Ctx: ctx}
	return response
}

func (r Response) Json(h H) {
	if h.Message == "" {
		h.Message = "Success"
	}
	r.Ctx.JSON(http.StatusOK, h)
}

func (r Response) Panic(statusCode int, h HErr) {
	if h.Message == "" {
		h.Message = "Request failed"
	}
	if h.ErrorCode == "" {
		h.ErrorCode = "4000"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(statusCode), h.ErrorCode, h.Message))
}
