package handler

import (
	"net/http"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	show, err := h.showRepo.FindByID(c.Request().Context(), showID)
	if err != nil {
		return err
	}

	counts, err := h.bookingRepo.AggregateSeatCounts(c.Request().Context(), show.ID)
	if err != nil {
		return err
	}

	resp := model.VerifyShowResponse{
		ShowID:    show.ID,
		SeatCount: show.SeatCount,
	}
	for _, sbc := range counts {
		resp.TotalBookings += sbc.Count
		switch {
		case sbc.Count == 1:
			resp.BookedSeats++
		case sbc.Count >= 2:
			resp.OversoldCount++
			resp.OversoldSeats = append(resp.OversoldSeats, model.OversoldSeat{
				SeatID:       sbc.SeatID,
				BookingCount: sbc.Count,
			})
		}
	}
	resp.UnbookedSeats = show.SeatCount - len(counts)
	resp.IsValid = resp.OversoldCount == 0

	if len(resp.OversoldSeats) > 0 {
		ids := make([]primitive.ObjectID, len(resp.OversoldSeats))
		for i, os := range resp.OversoldSeats {
			ids[i] = os.SeatID
		}

		seats, err := h.seatRepo.FindByIDs(c.Request().Context(), ids)
		if err != nil {
			return err
		}

		numberByID := make(map[primitive.ObjectID]int, len(seats))
		for _, s := range seats {
			numberByID[s.ID] = s.Number
		}

		for i := range resp.OversoldSeats {
			resp.OversoldSeats[i].Number = numberByID[resp.OversoldSeats[i].SeatID]
		}
	}

	return c.JSON(http.StatusOK, resp)
}
