package mongo

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type movieRepository struct {
	coll *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) domain.MovieRepository {
	return &movieRepository{
		coll: db.Collection("movies"),
	}
}

func (r *movieRepository) Create(ctx context.Context, movie *domain.Movie) error {
	_, err := r.coll.InsertOne(ctx, movie)
	return err
}

func (r *movieRepository) FindAll(ctx context.Context) ([]domain.Movie, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []domain.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *movieRepository) FindByID(ctx context.Context, id string) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &movie, err
}

func (r *movieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": movie.ID}, movie)
	return err
}

func (r *movieRepository) Delete(ctx context.Context, id string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
