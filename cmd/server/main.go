package main

import (
	"log"
	"os"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/database"
	"seojoonrp/ticket-rush-lab/internal/handler"
	appMiddleware "seojoonrp/ticket-rush-lab/internal/middleware"
	"seojoonrp/ticket-rush-lab/internal/repository"
	"seojoonrp/ticket-rush-lab/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	// database setup
	db, err := database.Connect(dbURI, dbName)
	if err != nil {
		panic(err)
	}
	defer database.Disconnect(db)

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

	// routes setup
	showRepo := repository.NewShowRepo(db)
	seatRepo := repository.NewSeatRepo(db)
	bookingRepo := repository.NewBookingRepo(db)

	showService := service.NewShowService(showRepo, seatRepo)
	bookingService := service.NewBookingService(showRepo, seatRepo, bookingRepo)

	handler := handler.NewHandler(showService, bookingService)

	e.POST("/shows", handler.RegisterShow)
	e.GET("/shows/:id/verify", handler.Verify)

	auth := e.Group("", appMiddleware.AuthMiddleware)
	auth.POST("/seats/:id/book", handler.Book)

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
