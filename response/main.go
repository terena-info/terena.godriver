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
	Panic(HErr)
	Status(int) Response
	SetHeader(string, string) Response
	SetCookie(string, string) Response
}

type Response struct {
	Ctx        *gin.Context
	statusCode int
}

func New(ctx *gin.Context) ResponseInterface {
	var response ResponseInterface = Response{Ctx: ctx, statusCode: http.StatusOK}
	return response
}

func (r Response) SetHeader(key string, value string) Response {
	r.Ctx.Header(key, value)
	return r
}

func (r Response) SetCookie(key string, value string) Response {
	r.Ctx.Header(key, value)
	return r
}

func (r Response) Status(statusCode int) Response {
	r.statusCode = statusCode
	return r
}

func (r Response) Json(h H) {
	if h.Message == "" {
		h.Message = "Success"
	}
	r.Ctx.JSON(http.StatusOK, h)
}

func (r Response) Panic(h HErr) {
	if h.Message == "" {
		h.Message = "Request failed"
	}
	if h.ErrorCode == "" {
		h.ErrorCode = "4000"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(r.statusCode), h.ErrorCode, h.Message))
}
