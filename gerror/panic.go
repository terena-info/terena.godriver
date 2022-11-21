package gerrors

import (
	"fmt"
	"strconv"
)

type E struct {
	Message   string
	ErrorCode string
}

func Panic(statusCode int, p E) {
	if p.ErrorCode == "" {
		p.ErrorCode = "5000"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(statusCode), p.ErrorCode, p.Message))
}

type Throw struct {
	Message   string
	ErrorCode string
}

func (e E) BadRequest() {
	if e.ErrorCode == "" {
		e.ErrorCode = "4000"
	}
	if e.Message == "" {
		e.Message = "Bad request"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.ErrorCode, e.Message))
}

func (e E) Unauthorized() {
	if e.ErrorCode == "" {
		e.ErrorCode = "4001"
	}
	if e.Message == "" {
		e.Message = "Unauthorized"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.ErrorCode, e.Message))
}

func (e E) Custom() {
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.ErrorCode, e.Message))
}

func (e E) InteralServerError() {
	if e.ErrorCode == "" {
		e.ErrorCode = "5000"
	}
	if e.Message == "" {
		e.Message = "Interal server error"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.ErrorCode, e.Message))
}
