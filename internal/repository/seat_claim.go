package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeatClaimRepo struct {
	rdb *redis.Client
}

func NewSeatClaimRepo(rdb *redis.Client) *SeatClaimRepo {
	return &SeatClaimRepo{rdb: rdb}
}

func (r *SeatClaimRepo) Claim(ctx context.Context, seatID primitive.ObjectID, userID string) (bool, error) {
	key := "seat:" + seatID.Hex()
	// SetNX(context, key, value, expiration)
	// value를 userID로 설정해 디버깅이 편하도록 (누가 이 좌석을 선점했는지 파악 가능)
	// expiration을 0으로 설정하면 만료되지 않음
	return r.rdb.SetNX(ctx, key, userID, 0).Result()
}
