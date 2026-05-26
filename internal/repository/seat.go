package repository

import (
	"context"
	"errors"
	"seojoonrp/ticket-rush-lab/internal/apperr"
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

func (r *SeatRepo) UpdateOnBook(ctx context.Context, id primitive.ObjectID, userID string) error {
	update := bson.M{"$set": bson.M{
		"status":  model.SeatOccupied,
		"user_id": userID,
	}}

	if _, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
		return err
	}

	return nil
}

func (r *SeatRepo) TryOccupy(ctx context.Context, id primitive.ObjectID, userID string) (bool, error) {
	filter := bson.M{
		"_id":    id,
		"status": model.SeatAvailable,
	}
	update := bson.M{"$set": bson.M{
		"status":  model.SeatOccupied,
		"user_id": userID,
	}}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return res.MatchedCount == 1, nil
}

func (r *SeatRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Seat, error) {
	var seat model.Seat
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&seat); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperr.ErrSeatNotFound
		}
		return nil, err
	}

	return &seat, nil
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
