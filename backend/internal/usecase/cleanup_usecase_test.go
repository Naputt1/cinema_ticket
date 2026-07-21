package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

func TestCleanup_ReleasesExpiredBookings(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	expiredAt := now.Add(-1 * time.Minute)

	expiredBookings := []domain.Booking{
		{
			ID:        "b-1",
			UserID:    "user-1",
			ShowtimeID: "st-1",
			SeatIDs:   []string{"seat-1", "seat-2"},
			Status:    domain.BookingStatusLocked,
			ExpiresAt: &expiredAt,
		},
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("FindExpiredLocked", ctx, mock.AnythingOfType("time.Time")).
		Run(func(args mock.Arguments) {
			passed := args.Get(1).(time.Time)
			assert.WithinDuration(t, time.Now(), passed, time.Second)
		}).
		Return(expiredBookings, nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, []string{"seat-1", "seat-2"}, domain.SeatStatusAvailable).Return(nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-1", "user-1").Return(true, nil)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-2", "user-1").Return(true, nil)

	var updatedBooking *domain.Booking
	mockBR.On("Update", ctx, mock.AnythingOfType("*domain.Booking")).Run(func(args mock.Arguments) {
		updatedBooking = args.Get(1).(*domain.Booking)
	}).Return(nil)

	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Return(nil).Times(2)

	mockAR := new(mockAuditRepo)
	mockAR.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil).Times(2)

	cleanup := NewBookingExpiryCleanup(mockBR, mockSeatR, mockLR, mockEP, mockAR)
	cleanup.cleanup(ctx)

	require.NotNil(t, updatedBooking)
	assert.Equal(t, domain.BookingStatusExpired, updatedBooking.Status)
	assert.Nil(t, updatedBooking.ExpiresAt)
	assert.WithinDuration(t, time.Now(), updatedBooking.UpdatedAt, time.Second)

	mockBR.AssertExpectations(t)
	mockSeatR.AssertExpectations(t)
	mockLR.AssertExpectations(t)
	mockEP.AssertExpectations(t)
	mockAR.AssertExpectations(t)
}

func TestCleanup_SkipsNoExpiredBookings(t *testing.T) {
	ctx := context.Background()

	mockBR := new(mockBookingRepo)
	mockBR.On("FindExpiredLocked", ctx, mock.AnythingOfType("time.Time")).Return([]domain.Booking{}, nil)

	cleanup := NewBookingExpiryCleanup(mockBR, new(mockSeatRepo), new(mockLockRepo), new(mockEventPub), new(mockAuditRepo))
	cleanup.cleanup(ctx)

	mockBR.AssertExpectations(t)
}

func TestCleanup_PublishesEventSeatReleasedAndBookingTimeout(t *testing.T) {
	ctx := context.Background()
	expiredAt := time.Now().Add(-1 * time.Minute)

	expiredBookings := []domain.Booking{
		{
			ID:        "b-1",
			UserID:    "user-1",
			ShowtimeID: "st-1",
			SeatIDs:   []string{"seat-1"},
			Status:    domain.BookingStatusLocked,
			ExpiresAt: &expiredAt,
		},
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("FindExpiredLocked", ctx, mock.AnythingOfType("time.Time")).Return(expiredBookings, nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, []string{"seat-1"}, domain.SeatStatusAvailable).Return(nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-1", "user-1").Return(true, nil)

	mockBR.On("Update", ctx, mock.AnythingOfType("*domain.Booking")).Return(nil)

	var publishedEvents []domain.Event
	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Run(func(args mock.Arguments) {
		publishedEvents = append(publishedEvents, args.Get(1).(domain.Event))
	}).Return(nil).Times(2)

	mockAR := new(mockAuditRepo)
	mockAR.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil).Times(2)

	cleanup := NewBookingExpiryCleanup(mockBR, mockSeatR, mockLR, mockEP, mockAR)
	cleanup.cleanup(ctx)

	require.Len(t, publishedEvents, 2)
	assert.Equal(t, domain.EventSeatReleased, publishedEvents[0].Type)
	assert.Equal(t, domain.EventBookingTimeout, publishedEvents[1].Type)

	mockEP.AssertExpectations(t)
	mockAR.AssertExpectations(t)
}
