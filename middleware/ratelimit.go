package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(50), int(50*3))

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}
		c.String(http.StatusOK, "rate limit!")
		c.Abort()
	}
}
