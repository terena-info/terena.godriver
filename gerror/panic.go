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
	M    string
	Code string
}

func (e Throw) BadRequest() {
	if e.Code == "" {
		e.Code = "4000"
	}
	if e.M == "" {
		e.M = "Bad request"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.Code, e.M))
}

func (e Throw) Unauthorized() {
	if e.Code == "" {
		e.Code = "4001"
	}
	if e.M == "" {
		e.M = "Unauthorized"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.Code, e.M))
}

func (e Throw) Custom() {
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.Code, e.M))
}

func (e Throw) InteralServerError() {
	if e.Code == "" {
		e.Code = "5000"
	}
	if e.M == "" {
		e.M = "Interal server error"
	}
	panic(fmt.Sprintf("%s-%s::%s", strconv.Itoa(400), e.Code, e.M))
}
