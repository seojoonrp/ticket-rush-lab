package repository

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatRepo struct {
	coll *mongo.Collection
}

func NewSeatRepo(db *mongo.Database) *SeatRepo {
	return &SeatRepo{coll: db.Collection("seats")}
}

func (r *SeatRepo) CreateMany(ctx context.Context, showID string, count int) ([]model.Seat, error) {
	shID, err := primitive.ObjectIDFromHex(showID)
	if err != nil {
		return nil, apperr.ErrInvalidID("show")
	}

	seats := make([]model.Seat, count)
	docs := make([]interface{}, count)
	for i := 0; i < count; i++ {
		seat := model.Seat{
			ID:     primitive.NewObjectID(),
			ShowID: shID,
			Number: i + 1,
			Status: model.SeatAvailable,
		}
		seats[i] = seat
		docs[i] = seat
	}

	if _, err := r.coll.InsertMany(ctx, docs); err != nil {
		return nil, err
	}

	return seats, nil
}
