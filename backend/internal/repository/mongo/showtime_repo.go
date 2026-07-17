package mongo

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type showtimeRepository struct {
	coll *mongo.Collection
}

func NewShowtimeRepository(db *mongo.Database) domain.ShowtimeRepository {
	return &showtimeRepository{
		coll: db.Collection("showtimes"),
	}
}

func (r *showtimeRepository) Create(ctx context.Context, showtime *domain.Showtime) error {
	_, err := r.coll.InsertOne(ctx, showtime)
	return err
}

func (r *showtimeRepository) CreateMany(ctx context.Context, showtimes []domain.Showtime) error {
	docs := make([]interface{}, len(showtimes))
	for i := range showtimes {
		docs[i] = showtimes[i]
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *showtimeRepository) FindByMovieID(ctx context.Context, movieID string) ([]domain.Showtime, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var showtimes []domain.Showtime
	if err := cursor.All(ctx, &showtimes); err != nil {
		return nil, err
	}
	return showtimes, nil
}

func (r *showtimeRepository) FindByID(ctx context.Context, id string) (*domain.Showtime, error) {
	var showtime domain.Showtime
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&showtime)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &showtime, err
}
