package usecase

import (
	"context"
	"fmt"

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
}

func NewAdminUsecase(
	bookingRepo domain.BookingRepository,
	showtimeRepo domain.ShowtimeRepository,
	movieRepo domain.MovieRepository,
	seatRepo domain.SeatRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
) *AdminUsecase {
	return &AdminUsecase{
		bookingRepo:  bookingRepo,
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		seatRepo:     seatRepo,
		userRepo:     userRepo,
		auditRepo:    auditRepo,
	}
}

type BookingFilter struct {
	MovieID string `form:"movie_id"`
	Date    string `form:"date"`
	UserID  string `form:"user_id"`
	Status  string `form:"status"`
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
		ID:        booking.ID,
		UserID:    booking.UserID,
		UserEmail: userEmail,
		Showtime:  *showtime,
		Movie:     *movie,
		Seats:     seats,
		Status:    booking.Status,
		ExpiresAt: booking.ExpiresAt,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}, nil
}

func (uc *AdminUsecase) GetAuditLogs(ctx context.Context, eventType string) ([]domain.AuditLog, error) {
	f := bson.M{}
	if eventType != "" {
		f["event_type"] = eventType
	}
	return uc.auditRepo.FindAll(ctx, f)
}
