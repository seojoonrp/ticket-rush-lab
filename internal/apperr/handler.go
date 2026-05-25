package apperr

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, map[string]string{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	// AppError가 아닌 시스템 에러
	c.Logger().Error(err)
	c.JSON(http.StatusInternalServerError, map[string]string{
		"code":  "INTERNAL",
		"error": "internal server error",
	})
}
