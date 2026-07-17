package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type MovieRepository interface {
	Create(ctx context.Context, movie *Movie) error
	FindAll(ctx context.Context) ([]Movie, error)
	FindByID(ctx context.Context, id string) (*Movie, error)
}

type ShowtimeRepository interface {
	Create(ctx context.Context, showtime *Showtime) error
	FindByMovieID(ctx context.Context, movieID string) ([]Showtime, error)
	FindByID(ctx context.Context, id string) (*Showtime, error)
	CreateMany(ctx context.Context, showtimes []Showtime) error
}

type SeatRepository interface {
	Create(ctx context.Context, seat *Seat) error
	BulkCreate(ctx context.Context, seats []Seat) error
	FindByShowtimeID(ctx context.Context, showtimeID string) ([]Seat, error)
	FindByIDs(ctx context.Context, ids []string) ([]Seat, error)
	FindByID(ctx context.Context, id string) (*Seat, error)
	UpdateStatus(ctx context.Context, id string, status SeatStatus) error
	BulkUpdateStatus(ctx context.Context, ids []string, status SeatStatus) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id string) (*Booking, error)
	FindByUserID(ctx context.Context, userID string) ([]Booking, error)
	FindAll(ctx context.Context, filter map[string]interface{}) ([]Booking, error)
	UpdateStatus(ctx context.Context, id string, status BookingStatus) error
	Update(ctx context.Context, booking *Booking) error
	FindExpiredLocked(ctx context.Context, now time.Time) ([]Booking, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	FindAll(ctx context.Context, filter map[string]interface{}) ([]AuditLog, error)
	BulkCreate(ctx context.Context, logs []AuditLog) error
}

type LockRepository interface {
	Acquire(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string, value string) (bool, error)
	Extend(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
}
