package main

import (
	"log"
	"net/http"
	"os"

	"github.com/eneridangelis/golangExpert/rate-limiter/limiter"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisClient := limiter.NewRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(limiter.RateLimiterMiddleware(redisClient))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
