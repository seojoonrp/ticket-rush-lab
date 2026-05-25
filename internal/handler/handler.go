package handler

import (
	"net/http"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"
	"seojoonrp/ticket-rush-lab/internal/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	showRepo       repository.ShowRepo
	seatRepo       repository.SeatRepo
	bookingRepo    repository.BookingRepo
	bookingService *service.BookingService
}

func NewHandler(
	shr repository.ShowRepo,
	ser repository.SeatRepo,
	br repository.BookingRepo,
	bs *service.BookingService,
) *Handler {
	return &Handler{
		showRepo:       shr,
		seatRepo:       ser,
		bookingRepo:    br,
		bookingService: bs,
	}
}

func (h *Handler) RegisterShow(c echo.Context) error {
	var req model.RegisterShowRequest
	if err := c.Bind(&req); err != nil {
		return apperr.ErrInvalidRequestBody
	}

	show, err := h.showRepo.Create(c.Request().Context(), req.SeatCount)
	if err != nil {
		return err
	}

	seats, err := h.seatRepo.CreateMany(c.Request().Context(), show.ID, req.SeatCount)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, model.RegisterShowResponse{
		Show:  *show,
		Seats: seats,
	})
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
