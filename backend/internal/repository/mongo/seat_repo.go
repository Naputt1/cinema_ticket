package mongo

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type seatRepository struct {
	coll *mongo.Collection
}

func NewSeatRepository(db *mongo.Database) domain.SeatRepository {
	return &seatRepository{
		coll: db.Collection("seats"),
	}
}

func (r *seatRepository) Create(ctx context.Context, seat *domain.Seat) error {
	_, err := r.coll.InsertOne(ctx, seat)
	return err
}

func (r *seatRepository) BulkCreate(ctx context.Context, seats []domain.Seat) error {
	docs := make([]interface{}, len(seats))
	for i := range seats {
		docs[i] = seats[i]
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *seatRepository) FindByShowtimeID(ctx context.Context, showtimeID string) ([]domain.Seat, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"showtime_id": showtimeID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var seats []domain.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) FindByIDs(ctx context.Context, ids []string) ([]domain.Seat, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var seats []domain.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) FindByID(ctx context.Context, id string) (*domain.Seat, error) {
	var seat domain.Seat
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&seat)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &seat, err
}

func (r *seatRepository) UpdateStatus(ctx context.Context, id string, status domain.SeatStatus) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *seatRepository) DeleteByShowtimeID(ctx context.Context, showtimeID string) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"showtime_id": showtimeID})
	return err
}

func (r *seatRepository) BulkUpdateStatus(ctx context.Context, ids []string, status domain.SeatStatus) error {
	_, err := r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"status": status}})
	return err
}
