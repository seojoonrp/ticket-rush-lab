package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/database"
	"seojoonrp/ticket-rush-lab/internal/handler"
	appMiddleware "seojoonrp/ticket-rush-lab/internal/middleware"
	"seojoonrp/ticket-rush-lab/internal/repository"
	"seojoonrp/ticket-rush-lab/internal/service"
	"strconv"
	"syscall"
	"time"

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

	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	redisAddr := os.Getenv("REDIS_ADDR")

	workerCount, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	bufferSize, _ := strconv.Atoi(os.Getenv("BUFFER_SIZE"))

	// database setup
	mdb, err := database.ConnectMongo(mongoURI, mongoDBName)
	if err != nil {
		panic(err)
	}
	defer database.DisconnectMongo(mdb)

	rdb, err := database.ConnectRedis(redisAddr)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

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

	showRepo := repository.NewShowRepo(mdb)
	seatRepo := repository.NewSeatRepo(mdb)
	bookingRepo := repository.NewBookingRepo(mdb)
	claimRepo := repository.NewSeatClaimRepo(rdb)

	pool := service.NewWorkerPool(seatRepo, bookingRepo, workerCount, bufferSize)

	showService := service.NewShowService(showRepo, seatRepo)
	bookingService := service.NewBookingService(showRepo, seatRepo, bookingRepo, claimRepo, pool)

	handler := handler.NewHandler(showService, bookingService, pool)

	// routes setup
	e.POST("/shows", handler.RegisterShow)
	e.GET("/shows/:id/verify", handler.Verify)
	e.GET("/health/queue", handler.QueueDepth)

	auth := e.Group("", appMiddleware.AuthMiddleware)
	auth.POST("/seats/:id/book", handler.Book)

	// start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	pool.Shutdown()
}
