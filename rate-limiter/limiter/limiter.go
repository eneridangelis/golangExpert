package limiter

import (
	"net/http"

	"github.com/eneridangelis/golangExpert/rate-limiter/config"

	"github.com/labstack/echo/v4"
)

type RateLimiter struct {
	strategy RepositoryStrategy
}

func NewRateLimiter(strategy RepositoryStrategy) *RateLimiter {
	return &RateLimiter{
		strategy: strategy,
	}
}

func (rl *RateLimiter) CheckLimit(ip, token string) (bool, error) {
	var key string
	var limit int

	if token != "" {
		key = getRateLimitKey("token", token)
		limit = config.GetEnvAsInt("DEFAULT_TOKEN_LIMIT", 20)
	} else {
		key = getRateLimitKey("ip", ip)
		limit = config.GetEnvAsInt("DEFAULT_IP_LIMIT", 10)
	}

	count, err := rl.strategy.IncrementAndGet(key, 1)
	if err != nil {
		return false, err
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}

func RateLimiterMiddleware(strategy RepositoryStrategy) echo.MiddlewareFunc {
	limiter := NewRateLimiter(strategy)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			token := c.Request().Header.Get("API_KEY")

			allowed, err := limiter.CheckLimit(ip, token)
			if err != nil {
				return err
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
				})
			}

			return next(c)
		}
	}
}

func getRateLimitKey(prefix, identifier string) string {
	return prefix + ":" + identifier
}
