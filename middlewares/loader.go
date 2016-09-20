package middlewares

import "github.com/gin-gonic/gin"

func Loader() gin.HandlerFunc {
	db := Database.Connect()

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}
