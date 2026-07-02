package handler

import (
	"net/http"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/middleware"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	showService    *service.ShowService
	bookingService *service.BookingService
	pool           *service.WorkerPool
}

func NewHandler(
	ss *service.ShowService,
	bs *service.BookingService,
	wp *service.WorkerPool,
) *Handler {
	return &Handler{
		showService:    ss,
		bookingService: bs,
		pool:           wp,
	}
}

func (h *Handler) RegisterShow(c echo.Context) error {
	var req model.RegisterShowRequest
	if err := c.Bind(&req); err != nil {
		return apperr.ErrInvalidRequestBody
	}

	resp, err := h.showService.RegisterShow(c.Request().Context(), req.SeatCount)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Book(c echo.Context) error {
	seatID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return apperr.ErrInvalidID("seat")
	}

	userID := middleware.GetUserID(c)

	if err := h.bookingService.Book(c.Request().Context(), seatID, userID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) QueueDepth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]int{"depth": h.pool.QueueDepth()})
}

func (h *Handler) Verify(c echo.Context) error {
	showID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return apperr.ErrInvalidID("show")
	}

	resp, err := h.bookingService.Verify(c.Request().Context(), showID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
