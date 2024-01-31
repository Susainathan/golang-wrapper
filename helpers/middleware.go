package helpers

import (
	"github.com/gin-gonic/gin"
)

func CommonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		defer func() {
			if err := recover(); err != nil {
				ApiException(c, ExceptionStruct{ExceptionKey: "Uncaught", Message: "Unexpected error occurred while processing the request. Please try after some time."}, 500)
			}
		}()

		c.Next()
	}
}
