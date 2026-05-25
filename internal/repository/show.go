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

type ShowRepo struct {
	coll *mongo.Collection
}

func NewShowRepo(db *mongo.Database) *ShowRepo {
	return &ShowRepo{coll: db.Collection("shows")}
}

func (r *ShowRepo) Create(ctx context.Context, seatCount int) (*model.Show, error) {
	show := model.Show{
		ID:        primitive.NewObjectID(),
		SeatCount: seatCount,
	}

	if _, err := r.coll.InsertOne(ctx, show); err != nil {
		return nil, err
	}

	return &show, nil
}

func (r *ShowRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Show, error) {
	var show model.Show
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&show); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperr.ErrShowNotFound
		}
		return nil, err
	}

	return &show, nil
}
