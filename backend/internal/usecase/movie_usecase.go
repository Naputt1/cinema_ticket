package usecase

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

type MovieUsecase struct {
	movieRepo    domain.MovieRepository
	showtimeRepo domain.ShowtimeRepository
	seatRepo     domain.SeatRepository
}

func NewMovieUsecase(
	movieRepo domain.MovieRepository,
	showtimeRepo domain.ShowtimeRepository,
	seatRepo domain.SeatRepository,
) *MovieUsecase {
	return &MovieUsecase{
		movieRepo:    movieRepo,
		showtimeRepo: showtimeRepo,
		seatRepo:     seatRepo,
	}
}

func (uc *MovieUsecase) GetMovies(ctx context.Context) ([]domain.Movie, error) {
	return uc.movieRepo.FindAll(ctx)
}

func (uc *MovieUsecase) GetMovie(ctx context.Context, id string) (*domain.Movie, error) {
	return uc.movieRepo.FindByID(ctx, id)
}

func (uc *MovieUsecase) GetShowtimes(ctx context.Context, movieID string) ([]domain.Showtime, error) {
	return uc.showtimeRepo.FindByMovieID(ctx, movieID)
}

func (uc *MovieUsecase) GetSeats(ctx context.Context, showtimeID string) ([]domain.Seat, error) {
	return uc.seatRepo.FindByShowtimeID(ctx, showtimeID)
}
