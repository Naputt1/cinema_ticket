package usecase

import (
	"context"
	"log"
	"time"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

type BookingExpiryCleanup struct {
	bookingRepo domain.BookingRepository
	seatRepo    domain.SeatRepository
	lockRepo    RedisLockRepo
	eventPub    domain.EventPublisher
	auditRepo   domain.AuditLogRepository
}

func NewBookingExpiryCleanup(
	bookingRepo domain.BookingRepository,
	seatRepo domain.SeatRepository,
	lockRepo RedisLockRepo,
	eventPub domain.EventPublisher,
	auditRepo domain.AuditLogRepository,
) *BookingExpiryCleanup {
	return &BookingExpiryCleanup{
		bookingRepo: bookingRepo,
		seatRepo:    seatRepo,
		lockRepo:    lockRepo,
		eventPub:    eventPub,
		auditRepo:   auditRepo,
	}
}

func (c *BookingExpiryCleanup) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *BookingExpiryCleanup) cleanup(ctx context.Context) {
	expired, err := c.bookingRepo.FindExpiredLocked(ctx, time.Now())
	if err != nil {
		log.Printf("failed to find expired bookings: %v", err)
		return
	}

	for _, booking := range expired {
		if err := c.seatRepo.BulkUpdateStatus(ctx, booking.SeatIDs, domain.SeatStatusAvailable); err != nil {
			log.Printf("failed to release seats for booking %s: %v", booking.ID, err)
			continue
		}

		for _, seatID := range booking.SeatIDs {
			lockKey := "lock:seat:" + booking.ShowtimeID + ":" + seatID
			c.lockRepo.Release(ctx, lockKey, booking.UserID)

		}

		booking.Status = domain.BookingStatusExpired
		booking.ExpiresAt = nil
		booking.UpdatedAt = time.Now()
		if err := c.bookingRepo.Update(ctx, &booking); err != nil {
			log.Printf("failed to update expired booking %s: %v", booking.ID, err)
		}

		c.eventPub.Publish(ctx, domain.Event{
			Type: domain.EventSeatReleased,
			Payload: map[string]interface{}{
				"booking_id":  booking.ID,
				"showtime_id": booking.ShowtimeID,
				"seat_ids":    booking.SeatIDs,
			},
		})

		c.eventPub.Publish(ctx, domain.Event{
			Type: domain.EventBookingTimeout,
			Payload: map[string]interface{}{
				"booking_id": booking.ID,
				"user_id":    booking.UserID,
			},
		})

		auditLog := &domain.AuditLog{
			ID:        booking.ID + "-timeout",
			EventType: domain.AuditEventBookingTimeout,
			UserID:    booking.UserID,
			Detail:    "booking " + booking.ID + " expired",
			CreatedAt: time.Now(),
		}
		c.auditRepo.Create(ctx, auditLog)

		auditLog2 := &domain.AuditLog{
			ID:        booking.ID + "-released",
			EventType: domain.AuditEventSeatReleased,
			UserID:    booking.UserID,
			Detail:    "seats released after booking timeout",
			CreatedAt: time.Now(),
		}
		c.auditRepo.Create(ctx, auditLog2)
	}
}
