package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminUsecase struct {
	bookingRepo  domain.BookingRepository
	showtimeRepo domain.ShowtimeRepository
	movieRepo    domain.MovieRepository
	seatRepo     domain.SeatRepository
	userRepo     domain.UserRepository
	auditRepo    domain.AuditLogRepository
	lockRepo     RedisLockRepo
}

func NewAdminUsecase(
	bookingRepo domain.BookingRepository,
	showtimeRepo domain.ShowtimeRepository,
	movieRepo domain.MovieRepository,
	seatRepo domain.SeatRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	lockRepo RedisLockRepo,
) *AdminUsecase {
	return &AdminUsecase{
		bookingRepo:  bookingRepo,
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		seatRepo:     seatRepo,
		userRepo:     userRepo,
		auditRepo:    auditRepo,
		lockRepo:     lockRepo,
	}
}

type BookingFilter struct {
	MovieID string `form:"movie_id"`
	Date    string `form:"date"`
	UserID  string `form:"user_id"`
	Status  string `form:"status"`
	Skip    int    `form:"skip"`
	Limit   int    `form:"limit"`
}

func (uc *AdminUsecase) GetBookings(ctx context.Context, filter BookingFilter) ([]domain.BookingDetail, error) {
	f := bson.M{}
	if filter.MovieID != "" {
		f["movie_id"] = filter.MovieID
	}
	if filter.UserID != "" {
		f["user_id"] = filter.UserID
	}
	if filter.Status != "" {
		f["status"] = filter.Status
	}
	if filter.Date != "" {
		f["created_at"] = bson.M{"$regex": filter.Date, "$options": "i"}
	}

	if filter.Skip > 0 {
		f["$skip"] = filter.Skip
	}
	if filter.Limit > 0 {
		f["$limit"] = filter.Limit
	}

	bookings, err := uc.bookingRepo.FindAll(ctx, f)
	if err != nil {
		return nil, err
	}

	details := make([]domain.BookingDetail, 0, len(bookings))
	for _, b := range bookings {
		d, err := uc.enrichBooking(ctx, &b)
		if err != nil {
			return nil, err
		}
		details = append(details, *d)
	}
	return details, nil
}

func (uc *AdminUsecase) enrichBooking(ctx context.Context, booking *domain.Booking) (*domain.BookingDetail, error) {
	showtime, err := uc.showtimeRepo.FindByID(ctx, booking.ShowtimeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find showtime: %w", err)
	}

	movie, err := uc.movieRepo.FindByID(ctx, showtime.MovieID)
	if err != nil {
		return nil, fmt.Errorf("failed to find movie: %w", err)
	}

	seats, err := uc.seatRepo.FindByIDs(ctx, booking.SeatIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find seats: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, booking.UserID)
	userEmail := ""
	if err == nil && user != nil {
		userEmail = user.Email
	}

	return &domain.BookingDetail{
		ID:          booking.ID,
		UserID:      booking.UserID,
		UserEmail:   userEmail,
		Showtime:    *showtime,
		Movie:       *movie,
		Seats:       seats,
		Status:      booking.Status,
		ExpiresAt:   booking.ExpiresAt,
		Price:       booking.Price,
		TotalAmount: booking.TotalAmount,
		CreatedAt:   booking.CreatedAt,
		UpdatedAt:   booking.UpdatedAt,
	}, nil
}

func (uc *AdminUsecase) GetAllMovies(ctx context.Context) ([]domain.Movie, error) {
	return uc.movieRepo.FindAll(ctx)
}

func (uc *AdminUsecase) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	movie.ID = uuid.New().String()
	movie.UpdatedAt = time.Now()
	return uc.movieRepo.Create(ctx, movie)
}

func (uc *AdminUsecase) UpdateMovie(ctx context.Context, movie *domain.Movie) error {
	movie.UpdatedAt = time.Now()
	return uc.movieRepo.Update(ctx, movie)
}

func (uc *AdminUsecase) DeleteMovie(ctx context.Context, id string) error {
	return uc.movieRepo.Delete(ctx, id)
}

func (uc *AdminUsecase) GetAllShowtimes(ctx context.Context) ([]domain.Showtime, error) {
	return uc.showtimeRepo.FindAll(ctx)
}

func (uc *AdminUsecase) CreateShowtime(ctx context.Context, movieID string, st *domain.Showtime) error {
	st.ID = uuid.New().String()
	st.MovieID = movieID
	st.UpdatedAt = time.Now()
	return uc.showtimeRepo.Create(ctx, st)
}

func (uc *AdminUsecase) UpdateShowtime(ctx context.Context, st *domain.Showtime) error {
	st.UpdatedAt = time.Now()
	return uc.showtimeRepo.Update(ctx, st)
}

func (uc *AdminUsecase) DeleteShowtime(ctx context.Context, id string) error {
	if err := uc.bookingRepo.DeleteByShowtimeID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete bookings: %w", err)
	}
	if err := uc.seatRepo.DeleteByShowtimeID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete seats: %w", err)
	}
	return uc.showtimeRepo.Delete(ctx, id)
}

func (uc *AdminUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.FindAll(ctx)
}

func (uc *AdminUsecase) UpdateUserRole(ctx context.Context, userID string, role domain.Role) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.Role = role
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

func (uc *AdminUsecase) GetAuditLogs(ctx context.Context, eventType string) ([]domain.AuditLog, error) {
	f := bson.M{}
	if eventType != "" {
		f["event_type"] = eventType
	}
	return uc.auditRepo.FindAll(ctx, f)
}

func (uc *AdminUsecase) ForceCancelBooking(ctx context.Context, bookingID string) error {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}
	if booking.Status != domain.BookingStatusLocked && booking.Status != domain.BookingStatusBooked {
		return fmt.Errorf("booking cannot be cancelled")
	}
	booking.Status = domain.BookingStatusCancelled
	booking.ExpiresAt = nil
	booking.UpdatedAt = time.Now()
	if err := uc.bookingRepo.Update(ctx, booking); err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	if err := uc.seatRepo.BulkUpdateStatus(ctx, booking.SeatIDs, domain.SeatStatusAvailable); err != nil {
		return fmt.Errorf("failed to release seats: %w", err)
	}
	for _, seatID := range booking.SeatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", booking.ShowtimeID, seatID)
		uc.lockRepo.Release(ctx, lockKey, "admin")
	}
	return nil
}
