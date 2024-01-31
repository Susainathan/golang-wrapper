package helpers

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetNewUUID() uuid.UUID {
	return uuid.New()
}

func ApiException(c *gin.Context, exception ExceptionStruct, statusCode int) {
	exception.Id = GetNewUUID()
	exception.Debug = string(debug.Stack())
	c.Set("apiExceptionMessage", exception.Message)
	c.AbortWithStatusJSON(statusCode, exception)
	return
}
