package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatOccupied  SeatStatus = "OCCUPIED"
)

type Show struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SeatCount int                `bson:"seat_count" json:"seatCount"`
}

type Seat struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	ShowID primitive.ObjectID `bson:"show_id" json:"showId"`
	Number int                `bson:"number" json:"number"`
	Status SeatStatus         `bson:"status" json:"status"`
	UserID string             `bson:"user_id" json:"userId"`
}

type Booking struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SeatID    primitive.ObjectID `bson:"seat_id" json:"seatId"`
	UserID    string             `bson:"user_id" json:"userId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

type RegisterShowRequest struct {
	SeatCount int `json:"seatCount"`
}

type RegisterShowResponse struct {
	Show  Show   `json:"show"`
	Seats []Seat `json:"seats"`
}

type OversoldSeat struct {
	SeatID       primitive.ObjectID `json:"seatId"`
	SeatNumber   int                `json:"seatNumber"`
	BookingCount int                `json:"bookingCount"`
}

type VerifyShowResponse struct {
	ShowID        primitive.ObjectID `json:"showId"`
	SeatCount     int                `json:"seatCount"`
	TotalBookings int                `json:"totalBookings"`

	UnbookedSeats int `json:"unbookedSeats"`
	BookedSeats   int `json:"bookedSeats"`
	OversoldCount int `json:"oversoldCount"`

	IsValid       bool           `json:"isValid"`
	OversoldSeats []OversoldSeat `json:"oversoldSeats"`
}
