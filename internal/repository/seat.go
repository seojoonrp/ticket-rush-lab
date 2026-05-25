package repository

import (
	"context"
	"seojoonrp/ticket-rush-lab/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatRepo struct {
	coll *mongo.Collection
}

func NewSeatRepo(db *mongo.Database) *SeatRepo {
	return &SeatRepo{coll: db.Collection("seats")}
}

func (r *SeatRepo) CreateMany(ctx context.Context, showID primitive.ObjectID, count int) ([]model.Seat, error) {
	seats := make([]model.Seat, count)
	docs := make([]interface{}, count)
	for i := 0; i < count; i++ {
		seat := model.Seat{
			ID:     primitive.NewObjectID(),
			ShowID: showID,
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

func (r *SeatRepo) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Seat, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var seats []model.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}

	return seats, nil
}
