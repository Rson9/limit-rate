package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rson9/limit-rate/internal/limiter"
)

func RateLimitMiddleware(l limiter.Limter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !l.Check() {
				return c.String(http.StatusTooManyRequests, "rate limit exceeded")
			}
			return next(c)
		}
	}
}
