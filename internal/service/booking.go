package service

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	showRepo    *repository.ShowRepo
	seatRepo    *repository.SeatRepo
	bookingRepo *repository.BookingRepo
	claimRepo   *repository.SeatClaimRepo
}

func NewBookingService(
	shr *repository.ShowRepo,
	ser *repository.SeatRepo,
	br *repository.BookingRepo,
	cr *repository.SeatClaimRepo,
) *BookingService {
	return &BookingService{
		showRepo:    shr,
		seatRepo:    ser,
		bookingRepo: br,
		claimRepo:   cr,
	}
}

func (s *BookingService) Book(ctx context.Context, seatID primitive.ObjectID, userID string) error {
	claimed, err := s.claimRepo.Claim(ctx, seatID, userID)
	if err != nil {
		return err
	}
	if !claimed {
		return apperr.ErrSeatTaken
	}

	seat, err := s.seatRepo.FindByID(ctx, seatID)
	if err != nil {
		return err
	}

	if err := s.seatRepo.UpdateOnBook(ctx, seatID, userID); err != nil {
		return err
	}

	_, err = s.bookingRepo.Create(ctx, seat.ShowID, seatID, userID)
	if err != nil {
		return err
	}

	return nil
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
