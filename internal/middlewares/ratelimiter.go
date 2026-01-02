package middlewares

import (
	"net/http"
	"time"

	"github.com/as-ifn-at/REST/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Token bucket
func RateLimit() gin.HandlerFunc {

	limiter := rate.NewLimiter(per(time.Second, common.MaxNoOfRequestAllowed), common.MaxNoOfRequestAllowed)
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}
		ctx.Next()
	}
}

func per(duration time.Duration, eventCount int) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
