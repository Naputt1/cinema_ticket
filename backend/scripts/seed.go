package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	mongoRepo "github.com/naputt/cinema-ticket-backend/internal/repository/mongo"
	"github.com/naputt/cinema-ticket-backend/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())
	mongoDB := mongoClient.Database(cfg.Mongo.DB)

	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	redisClient.FlushDB(ctx)

	for _, coll := range []string{"users", "movies", "showtimes", "seats", "bookings", "audit_logs"} {
		mongoDB.Collection(coll).DeleteMany(ctx, bson.M{})
	}

	movieRepo := mongoRepo.NewMovieRepository(mongoDB)
	showtimeRepo := mongoRepo.NewShowtimeRepository(mongoDB)
	seatRepo := mongoRepo.NewSeatRepository(mongoDB)
	userRepo := mongoRepo.NewUserRepository(mongoDB)

	movies := []domain.Movie{
		{ID: uuid.New().String(), Title: "Inception", Description: "A thief who steals corporate secrets through dream-sharing technology is given the task of planting an idea into the mind of a C.E.O.", Duration: 148, PosterURL: "https://placehold.co/300x450?text=Inception", Genre: "Sci-Fi"},
		{ID: uuid.New().String(), Title: "The Matrix", Description: "A computer hacker learns about the true nature of reality and his role in the war against its controllers.", Duration: 136, PosterURL: "https://placehold.co/300x450?text=Matrix", Genre: "Action"},
		{ID: uuid.New().String(), Title: "Interstellar", Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.", Duration: 169, PosterURL: "https://placehold.co/300x450?text=Interstellar", Genre: "Sci-Fi"},
	}

	for i := range movies {
		movieRepo.Create(ctx, &movies[i])
	}
	log.Printf("seeded %d movies", len(movies))

	now := time.Now().Truncate(time.Hour)
	screens := []string{"Screen 1", "Screen 2"}
	var showtimes []domain.Showtime

	for _, movie := range movies {
		for _, screen := range screens {
			for h := 10; h <= 20; h += 5 {
				st := now.Add(time.Duration(h) * time.Hour)
				et := st.Add(time.Duration(movie.Duration) * time.Minute)
				showtimes = append(showtimes, domain.Showtime{
					ID:        uuid.New().String(),
					MovieID:   movie.ID,
					Screen:    screen,
					StartTime: st,
					EndTime:   et,
					Price:     12.0 + float64(h%3)*3.0,
				})
			}
		}
	}
	showtimeRepo.CreateMany(ctx, showtimes)
	log.Printf("seeded %d showtimes", len(showtimes))

	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, showtime := range showtimes {
		var seats []domain.Seat
		for _, row := range rows {
			for num := 1; num <= 8; num++ {
				seats = append(seats, domain.Seat{
					ID:         uuid.New().String(),
					ShowtimeID: showtime.ID,
					Row:        row,
					Number:     num,
					Status:     domain.SeatStatusAvailable,
					Screen:     showtime.Screen,
				})
			}
		}
		seatRepo.BulkCreate(ctx, seats)
	}
	log.Printf("seeded seats for %d showtimes", len(showtimes))

	adminUser := &domain.User{
		ID:        uuid.New().String(),
		Email:     "admin@example.com",
		Name:      "Admin",
		AvatarURL: "https://placehold.co/100?text=Admin",
		Role:      domain.RoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(ctx, adminUser)
	log.Printf("seeded admin user: admin@example.com")

	log.Println("seed complete")
}
