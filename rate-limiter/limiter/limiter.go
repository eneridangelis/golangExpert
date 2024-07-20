package limiter

import (
	"net/http"

	"github.com/eneridangelis/golangExpert/rate-limiter/config"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func RateLimiterMiddleware(client *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			token := c.Request().Header.Get("API_KEY")

			var key string
			var limit int

			if token != "" {
				key = GetRateLimitKey("token", token)
				limit = config.GetEnvAsInt("DEFAULT_TOKEN_LIMIT", 20)
			} else {
				key = GetRateLimitKey("ip", ip)
				limit = config.GetEnvAsInt("DEFAULT_IP_LIMIT", 10)
			}

			count, err := IncrementRequestCount(client, key, 1)
			if err != nil {
				return err
			}

			if count > int64(limit) {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
				})
			}

			return next(c)
		}
	}
}
