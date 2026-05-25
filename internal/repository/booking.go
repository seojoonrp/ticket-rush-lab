package repository

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/apperr"
	"seojoonrp/ticket-rush-lab/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo struct {
	coll *mongo.Collection
}

func NewBookingRepo(db *mongo.Database) *BookingRepo {
	return &BookingRepo{coll: db.Collection("bookings")}
}

func (r *BookingRepo) Create(ctx context.Context, seatID string, userID string) (*model.Booking, error) {
	seID, err := primitive.ObjectIDFromHex(seatID)
	if err != nil {
		return nil, apperr.ErrInvalidID("seat")
	}

	booking := model.Booking{
		ID:        primitive.NewObjectID(),
		SeatID:    seID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if _, err := r.coll.InsertOne(ctx, booking); err != nil {
		return nil, err
	}

	return &booking, nil
}
