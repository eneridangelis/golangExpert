package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eneridangelis/golangExpert/rate-limiter/limiter"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*echo.Echo, *redis.Client) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	e := echo.New()
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
	e.Use(limiter.RateLimiterMiddleware(client))

	client.FlushAll(client.Context())

	return e, client
}

func TestRateLimiterWithinLimit(t *testing.T) {
	e, _ := setupTest(t)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	e.GET("/", handler)

	limit := 10
	for i := 0; i < limit; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestRateLimiterExceedsLimit(t *testing.T) {
	e, _ := setupTest(t)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	e.GET("/", handler)

	limit := 10
	for i := 0; i < limit; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRateLimiterWithTokenExceedsLimit(t *testing.T) {
	e, _ := setupTest(t)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	e.GET("/", handler)

	token := "abc123"
	limit := 20
	for i := 0; i < limit; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("API_KEY", token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}
