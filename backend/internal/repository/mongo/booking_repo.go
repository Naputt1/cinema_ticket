package mongo

import (
	"context"
	"time"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookingRepository struct {
	coll *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) domain.BookingRepository {
	return &bookingRepository{
		coll: db.Collection("bookings"),
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	_, err := r.coll.InsertOne(ctx, booking)
	return err
}

func (r *bookingRepository) FindByID(ctx context.Context, id string) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &booking, err
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Booking, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []domain.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindAll(ctx context.Context, filter map[string]interface{}) ([]domain.Booking, error) {
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []domain.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	return err
}

func (r *bookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	booking.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": booking.ID}, booking)
	return err
}

func (r *bookingRepository) FindExpiredLocked(ctx context.Context, now time.Time) ([]domain.Booking, error) {
	cursor, err := r.coll.Find(ctx, bson.M{
		"status":     domain.BookingStatusLocked,
		"expires_at": bson.M{"$lte": now},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []domain.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
