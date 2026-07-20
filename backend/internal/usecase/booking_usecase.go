package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

type RedisLockRepo interface {
	Acquire(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string, value string) (bool, error)
	Extend(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
}

type BookingUsecase struct {
	bookingRepo  domain.BookingRepository
	seatRepo     domain.SeatRepository
	showtimeRepo domain.ShowtimeRepository
	movieRepo    domain.MovieRepository
	lockRepo     RedisLockRepo
	eventPub     domain.EventPublisher
	auditRepo    domain.AuditLogRepository
}

func NewBookingUsecase(
	bookingRepo domain.BookingRepository,
	seatRepo domain.SeatRepository,
	showtimeRepo domain.ShowtimeRepository,
	movieRepo domain.MovieRepository,
	lockRepo RedisLockRepo,
	eventPub domain.EventPublisher,
	auditRepo domain.AuditLogRepository,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo:  bookingRepo,
		seatRepo:     seatRepo,
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		lockRepo:     lockRepo,
		eventPub:     eventPub,
		auditRepo:    auditRepo,
	}
}

func (uc *BookingUsecase) LockSeat(ctx context.Context, userID string, showtimeID string, seatIDs []string) (*domain.Booking, error) {
	showtime, err := uc.showtimeRepo.FindByID(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find showtime: %w", err)
	}
	if showtime == nil {
		return nil, fmt.Errorf("showtime not found")
	}

	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", showtimeID, seatID)
		acquired, err := uc.lockRepo.Acquire(ctx, lockKey, userID, 5*time.Minute)
		if err != nil {
			uc.publishAudit(ctx, domain.AuditEventSystemError, userID, fmt.Sprintf("lock failed for seat %s: %v", seatID, err))
			uc.rollbackLocks(ctx, userID, showtimeID, seatIDs)
			return nil, fmt.Errorf("lock system error for seat %s", seatID)
		}
		if !acquired {
			uc.rollbackLocks(ctx, userID, showtimeID, seatIDs)
			return nil, fmt.Errorf("seat %s is already locked or booked", seatID)
		}
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	booking := &domain.Booking{
		ID:          uuid.New().String(),
		UserID:      userID,
		ShowtimeID:  showtimeID,
		SeatIDs:     seatIDs,
		Status:      domain.BookingStatusLocked,
		ExpiresAt:   &expiresAt,
		Price:       showtime.Price,
		TotalAmount: showtime.Price * float64(len(seatIDs)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.bookingRepo.Create(ctx, booking); err != nil {
		uc.rollbackLocks(ctx, userID, showtimeID, seatIDs)
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	if err := uc.seatRepo.BulkUpdateStatus(ctx, seatIDs, domain.SeatStatusLocked); err != nil {
		log.Printf("warning: failed to update seat status to locked: %v", err)
	}

	uc.eventPub.Publish(ctx, domain.Event{
		Type: domain.EventSeatLocked,
		Payload: map[string]interface{}{
			"booking_id":  booking.ID,
			"showtime_id": showtimeID,
			"seat_ids":    seatIDs,
			"user_id":     userID,
			"expires_at":  expiresAt,
		},
	})

	return booking, nil
}

func (uc *BookingUsecase) ExtendLock(ctx context.Context, userID string, bookingID string) error {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}
	if booking.UserID != userID {
		return fmt.Errorf("booking does not belong to user")
	}
	if booking.Status != domain.BookingStatusLocked {
		return fmt.Errorf("booking is not in locked state")
	}

	newExpiry := time.Now().Add(5 * time.Minute)
	for _, seatID := range booking.SeatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", booking.ShowtimeID, seatID)
		if _, err := uc.lockRepo.Extend(ctx, lockKey, userID, 5*time.Minute); err != nil {
			return fmt.Errorf("failed to extend lock for seat %s: %w", seatID, err)
		}
	}

	booking.ExpiresAt = &newExpiry
	booking.UpdatedAt = time.Now()
	return uc.bookingRepo.Update(ctx, booking)
}

func (uc *BookingUsecase) PayBooking(ctx context.Context, userID string, bookingID string) (*domain.Booking, error) {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("booking does not belong to user")
	}
	if booking.Status != domain.BookingStatusLocked {
		return nil, fmt.Errorf("booking is not in locked state")
	}

	booking.Status = domain.BookingStatusBooked
	booking.ExpiresAt = nil
	booking.UpdatedAt = time.Now()

	if err := uc.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	if err := uc.seatRepo.BulkUpdateStatus(ctx, booking.SeatIDs, domain.SeatStatusBooked); err != nil {
		log.Printf("warning: failed to update seat status to booked: %v", err)
	}

	for _, seatID := range booking.SeatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", booking.ShowtimeID, seatID)
		uc.lockRepo.Release(ctx, lockKey, userID)
	}

	uc.eventPub.Publish(ctx, domain.Event{
		Type: domain.EventSeatBooked,
		Payload: map[string]interface{}{
			"booking_id":  booking.ID,
			"showtime_id": booking.ShowtimeID,
			"seat_ids":    booking.SeatIDs,
			"user_id":     userID,
		},
	})

	uc.eventPub.Publish(ctx, domain.Event{
		Type: domain.EventBookingSuccess,
		Payload: map[string]interface{}{
			"booking_id": booking.ID,
			"user_id":    userID,
			"seat_ids":   booking.SeatIDs,
		},
	})

	uc.publishAudit(ctx, domain.AuditEventBookingSuccess, userID, fmt.Sprintf("booking %s completed for seats %v", booking.ID, booking.SeatIDs))

	return booking, nil
}

func (uc *BookingUsecase) CancelBooking(ctx context.Context, userID string, bookingID string) error {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}
	if booking.UserID != userID {
		return fmt.Errorf("booking does not belong to user")
	}
	if booking.Status != domain.BookingStatusLocked {
		return fmt.Errorf("booking is not in locked state")
	}

	return uc.releaseBooking(ctx, booking, domain.BookingStatusCancelled, domain.AuditEventSeatReleased)
}

func (uc *BookingUsecase) RefundBooking(ctx context.Context, userID string, bookingID string) error {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to find booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}
	if booking.UserID != userID {
		return fmt.Errorf("booking does not belong to user")
	}
	if booking.Status != domain.BookingStatusLocked && booking.Status != domain.BookingStatusBooked {
		return fmt.Errorf("booking cannot be cancelled")
	}

	return uc.releaseBooking(ctx, booking, domain.BookingStatusCancelled, domain.AuditEventSeatReleased)
}

func (uc *BookingUsecase) releaseBooking(ctx context.Context, booking *domain.Booking, newStatus domain.BookingStatus, auditEvent domain.AuditEventType) error {
	booking.Status = newStatus
	booking.ExpiresAt = nil
	booking.UpdatedAt = time.Now()

	if err := uc.bookingRepo.Update(ctx, booking); err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	if err := uc.seatRepo.BulkUpdateStatus(ctx, booking.SeatIDs, domain.SeatStatusAvailable); err != nil {
		log.Printf("warning: failed to update seat status to available: %v", err)
	}

	for _, seatID := range booking.SeatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", booking.ShowtimeID, seatID)
		uc.lockRepo.Release(ctx, lockKey, booking.UserID)
	}

	uc.eventPub.Publish(ctx, domain.Event{
		Type: domain.EventSeatReleased,
		Payload: map[string]interface{}{
			"booking_id":  booking.ID,
			"showtime_id": booking.ShowtimeID,
			"seat_ids":    booking.SeatIDs,
		},
	})

	uc.publishAudit(ctx, auditEvent, booking.UserID, fmt.Sprintf("seats %v released from booking %s", booking.SeatIDs, booking.ID))

	return nil
}

func (uc *BookingUsecase) GetUserBookings(ctx context.Context, userID string) ([]domain.BookingDetail, error) {
	bookings, err := uc.bookingRepo.FindByUserID(ctx, userID)
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

func (uc *BookingUsecase) GetBooking(ctx context.Context, bookingID string) (*domain.BookingDetail, error) {
	booking, err := uc.bookingRepo.FindByID(ctx, bookingID)
	if err != nil || booking == nil {
		return nil, err
	}
	return uc.enrichBooking(ctx, booking)
}

func (uc *BookingUsecase) enrichBooking(ctx context.Context, booking *domain.Booking) (*domain.BookingDetail, error) {
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

	return &domain.BookingDetail{
		ID:          booking.ID,
		UserID:      booking.UserID,
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

func (uc *BookingUsecase) rollbackLocks(ctx context.Context, userID string, showtimeID string, seatIDs []string) {
	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", showtimeID, seatID)
		uc.lockRepo.Release(ctx, lockKey, userID)
	}
}

func (uc *BookingUsecase) publishAudit(ctx context.Context, eventType domain.AuditEventType, userID string, detail string) {
	auditLog := &domain.AuditLog{
		ID:        uuid.New().String(),
		EventType: eventType,
		UserID:    userID,
		Detail:    detail,
		CreatedAt: time.Now(),
	}
	if err := uc.auditRepo.Create(ctx, auditLog); err != nil {
		log.Printf("warning: failed to create audit log: %v", err)
	}
}
