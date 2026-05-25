package main

import (
	"log"
	"seojoonrp/ticket-rush-lab/internal/apperr"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// echo setup
	e := echo.New()
	e.HTTPErrorHandler = apperr.Handler
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Printf("%s %s %d %s ERROR: %v\n",
					v.Method, v.URI, v.Status, v.Latency, v.Error)
			} else {
				log.Printf("%s %s %d %s\n",
					v.Method, v.URI, v.Status, v.Latency)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
