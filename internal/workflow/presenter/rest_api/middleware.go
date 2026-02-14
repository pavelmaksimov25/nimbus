package restapi

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewCorsMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{"*"}
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func NewSecureMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		ContentSecurityPolicy: "default-src 'self'",
	})
}

func NewRateLimiterMiddleware() gin.HandlerFunc {
	rateStr := "100-M"
	if r := os.Getenv("RATE_LIMIT"); r != "" {
		rateStr = r
	}

	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		rate = limiter.Rate{
			Period: time.Minute,
			Limit:  100,
		}
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance)
}
