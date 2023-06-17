package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"senao-auth-srv/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, e := range c.Errors {
			err := e.Err
			if customError, ok := err.(*errors.CustomError); ok {
				c.JSON(customError.Code, gin.H{
					"success": customError.Success,
					"reason":  customError.Reason,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": customError.Success,
					"reason":  customError.Reason,
				})
			}
			return
		}
	}
}
