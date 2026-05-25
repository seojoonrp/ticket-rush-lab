package middleware

import (
	"seojoonrp/ticket-rush-lab/internal/apperr"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Request().Header.Get("X-User-ID")
		if uid == "" {
			return apperr.ErrUnauthorized
		}

		c.Set("user_id", uid)

		return next(c)
	}
}

func GetUserID(c echo.Context) string {
	uid, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return uid
}
