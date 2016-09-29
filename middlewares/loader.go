package middlewares

import "github.com/gin-gonic/gin"

func Loader() gin.HandlerFunc {
	return func(c *gin.Context) {
		DB := Database.Connect("production")
		// Initialize Database COnnection
		c.Set("DB", DB)
		c.Next()
	}
}
