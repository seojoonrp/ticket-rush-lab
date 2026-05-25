package service

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"seojoonrp/ticket-rush-lab/internal/repository"
)

type ShowService struct {
	showRepo *repository.ShowRepo
	seatRepo *repository.SeatRepo
}

func NewShowService(
	shr *repository.ShowRepo,
	ser *repository.SeatRepo,
) *ShowService {
	return &ShowService{
		showRepo: shr,
		seatRepo: ser,
	}
}

func (s *ShowService) RegisterShow(ctx context.Context, seatCount int) (*model.RegisterShowResponse, error) {
	if seatCount <= 0 {
		return nil, apperr.ErrInvalidSeatCount
	}

	show, err := s.showRepo.Create(ctx, seatCount)
	if err != nil {
		return nil, err
	}

	seats, err := s.seatRepo.CreateMany(ctx, show.ID, seatCount)
	if err != nil {
		return nil, err
	}

	return &model.RegisterShowResponse{
		Show:  *show,
		Seats: seats,
	}, nil
}
