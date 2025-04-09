package logger

import (
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context, e error) {
	_, file, line, ok := runtime.Caller(1)
	requestID := c.GetHeader("X-Request-ID")
	if ok {
		log.Println(requestID, file, line, e)
	} else {
		log.Println(requestID, e)
	}
}
