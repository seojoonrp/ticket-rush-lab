package service

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	showRepo    *repository.ShowRepo
	seatRepo    *repository.SeatRepo
	bookingRepo *repository.BookingRepo
}

func NewBookingService(
	shr *repository.ShowRepo,
	ser *repository.SeatRepo,
	br *repository.BookingRepo,
) *BookingService {
	return &BookingService{
		showRepo:    shr,
		seatRepo:    ser,
		bookingRepo: br,
	}
}

func (s *BookingService) Verify(ctx context.Context, showID primitive.ObjectID) (*model.VerifyShowResponse, error) {
	show, err := s.showRepo.FindByID(ctx, showID)
	if err != nil {
		return nil, err
	}

	counts, err := s.bookingRepo.AggregateSeatCounts(ctx, show.ID)
	if err != nil {
		return nil, err
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

		seats, err := s.seatRepo.FindByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}

		numberByID := make(map[primitive.ObjectID]int, len(seats))
		for _, st := range seats {
			numberByID[st.ID] = st.Number
		}

		for i := range resp.OversoldSeats {
			resp.OversoldSeats[i].Number = numberByID[resp.OversoldSeats[i].SeatID]
		}
	}

	return &resp, nil
}
