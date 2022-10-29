package errors

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
