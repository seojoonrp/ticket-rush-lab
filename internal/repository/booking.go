package repository

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo struct {
	coll *mongo.Collection
}

func NewBookingRepo(db *mongo.Database) *BookingRepo {
	return &BookingRepo{coll: db.Collection("bookings")}
}

func (r *BookingRepo) Create(ctx context.Context, showID, seatID primitive.ObjectID, userID string) (*model.Booking, error) {
	booking := model.Booking{
		ID:        primitive.NewObjectID(),
		ShowID:    showID,
		SeatID:    seatID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if _, err := r.coll.InsertOne(ctx, booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepo) AggregateSeatCounts(ctx context.Context, showID primitive.ObjectID) ([]model.SeatBookingCount, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "show_id", Value: showID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$seat_id"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.SeatBookingCount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
