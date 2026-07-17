package domain

import "time"

type BookingStatus string

const (
	BookingStatusLocked   BookingStatus = "LOCKED"
	BookingStatusBooked   BookingStatus = "BOOKED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusExpired  BookingStatus = "EXPIRED"
)

type Booking struct {
	ID         string        `json:"id" bson:"_id"`
	UserID     string        `json:"user_id" bson:"user_id"`
	ShowtimeID string        `json:"showtime_id" bson:"showtime_id"`
	SeatIDs    []string      `json:"seat_ids" bson:"seat_ids"`
	Status     BookingStatus `json:"status" bson:"status"`
	ExpiresAt  *time.Time    `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusLocked    SeatStatus = "LOCKED"
	SeatStatusBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID         string     `json:"id" bson:"_id"`
	ShowtimeID string     `json:"showtime_id" bson:"showtime_id"`
	Row        string     `json:"row" bson:"row"`
	Number     int        `json:"number" bson:"number"`
	Status     SeatStatus `json:"status" bson:"status"`
	Screen     string     `json:"screen" bson:"screen"`
}
