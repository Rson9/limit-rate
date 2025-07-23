package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rson9/limit-rate/internal/di"
	limterMiddleware "github.com/rson9/limit-rate/internal/middleware"
)

func main() {
	e := echo.New()

	l, _ := di.InitLimiter()
	e.Use(limterMiddleware.RateLimitMiddleware(l))
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
