package handler

import (
	"net/http"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	showRepo    repository.ShowRepo
	seatRepo    repository.SeatRepo
	bookingRepo repository.BookingRepo
}

func NewHandler(shr repository.ShowRepo, ser repository.SeatRepo, br repository.BookingRepo) *Handler {
	return &Handler{
		showRepo:    shr,
		seatRepo:    ser,
		bookingRepo: br,
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

	seats, err := h.seatRepo.CreateMany(c.Request().Context(), show.ID.Hex(), req.SeatCount)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, model.RegisterShowResponse{
		Show:  *show,
		Seats: seats,
	})
}
