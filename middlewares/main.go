package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorRecovery(c *gin.Context, err interface{}) {
	if err, ok := err.(string); ok {
		errorCode := "5000"
		errorHttpStatusCode := http.StatusInternalServerError
		errMessage := err

		errSplited := strings.Split(err, "::")
		if len(errSplited) > 1 {
			errMessage = errSplited[1]
			splitErrorCode := strings.Split(errSplited[0], "-")
			if len(splitErrorCode) > 1 {
				errorHttpStatusCode, _ = strconv.Atoi(splitErrorCode[0])
				errorCode = splitErrorCode[1]
				c.JSON(errorHttpStatusCode, gin.H{"message": errMessage, "errorCode": errorCode})
				c.Abort()
				return
			}
			errorHttpStatusCode, _ = strconv.Atoi(errSplited[0])
		}

		c.JSON(errorHttpStatusCode, gin.H{"message": errMessage, "errorCode": errorCode})
	}
	c.Abort()
}
