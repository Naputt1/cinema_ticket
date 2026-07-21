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

type mockBookingRepo struct{ mock.Mock }
type mockSeatRepo struct{ mock.Mock }
type mockShowtimeRepo struct{ mock.Mock }
type mockMovieRepo struct{ mock.Mock }
type mockLockRepo struct{ mock.Mock }
type mockEventPub struct{ mock.Mock }
type mockAuditRepo struct{ mock.Mock }

func (m *mockBookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockBookingRepo) FindByID(ctx context.Context, id string) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}
func (m *mockBookingRepo) FindByUserID(ctx context.Context, uid string) ([]domain.Booking, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *mockBookingRepo) FindAll(ctx context.Context, f map[string]interface{}) ([]domain.Booking, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *mockBookingRepo) UpdateStatus(ctx context.Context, id string, s domain.BookingStatus) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *mockBookingRepo) Update(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockBookingRepo) FindExpiredLocked(ctx context.Context, now time.Time) ([]domain.Booking, error) {
	args := m.Called(ctx, now)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *mockBookingRepo) DeleteByShowtimeID(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockSeatRepo) Create(ctx context.Context, s *domain.Seat) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockSeatRepo) BulkCreate(ctx context.Context, s []domain.Seat) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockSeatRepo) FindByShowtimeID(ctx context.Context, id string) ([]domain.Seat, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]domain.Seat), args.Error(1)
}
func (m *mockSeatRepo) FindByIDs(ctx context.Context, ids []string) ([]domain.Seat, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]domain.Seat), args.Error(1)
}
func (m *mockSeatRepo) FindByID(ctx context.Context, id string) (*domain.Seat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Seat), args.Error(1)
}
func (m *mockSeatRepo) UpdateStatus(ctx context.Context, id string, s domain.SeatStatus) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *mockSeatRepo) BulkUpdateStatus(ctx context.Context, ids []string, s domain.SeatStatus) error {
	return m.Called(ctx, ids, s).Error(0)
}
func (m *mockSeatRepo) DeleteByShowtimeID(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockShowtimeRepo) Create(ctx context.Context, s *domain.Showtime) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockShowtimeRepo) FindAll(ctx context.Context) ([]domain.Showtime, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Showtime), args.Error(1)
}
func (m *mockShowtimeRepo) FindByMovieID(ctx context.Context, id string) ([]domain.Showtime, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]domain.Showtime), args.Error(1)
}
func (m *mockShowtimeRepo) FindByID(ctx context.Context, id string) (*domain.Showtime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Showtime), args.Error(1)
}
func (m *mockShowtimeRepo) CreateMany(ctx context.Context, s []domain.Showtime) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockShowtimeRepo) Update(ctx context.Context, s *domain.Showtime) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockShowtimeRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockMovieRepo) Create(ctx context.Context, mv *domain.Movie) error {
	return m.Called(ctx, mv).Error(0)
}
func (m *mockMovieRepo) FindAll(ctx context.Context) ([]domain.Movie, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Movie), args.Error(1)
}
func (m *mockMovieRepo) FindByID(ctx context.Context, id string) (*domain.Movie, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Movie), args.Error(1)
}
func (m *mockMovieRepo) Update(ctx context.Context, mv *domain.Movie) error {
	return m.Called(ctx, mv).Error(0)
}
func (m *mockMovieRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockLockRepo) Acquire(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, ttl)
	return args.Bool(0), args.Error(1)
}
func (m *mockLockRepo) Release(ctx context.Context, key, value string) (bool, error) {
	args := m.Called(ctx, key, value)
	return args.Bool(0), args.Error(1)
}
func (m *mockLockRepo) Extend(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *mockEventPub) Publish(ctx context.Context, e domain.Event) error {
	return m.Called(ctx, e).Error(0)
}

func (m *mockAuditRepo) Create(ctx context.Context, l *domain.AuditLog) error {
	return m.Called(ctx, l).Error(0)
}
func (m *mockAuditRepo) FindAll(ctx context.Context, f map[string]interface{}) ([]domain.AuditLog, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}
func (m *mockAuditRepo) BulkCreate(ctx context.Context, l []domain.AuditLog) error {
	return m.Called(ctx, l).Error(0)
}

func TestLockSeat_SetsExpiresAtToNowPlus5m(t *testing.T) {
	var (
		ctx        = context.Background()
		userID     = "user-1"
		showtimeID = "showtime-1"
		seatIDs    = []string{"seat-1", "seat-2"}
	)

	showtime := &domain.Showtime{
		ID:    showtimeID,
		Price: 12.50,
	}

	mockSR := new(mockShowtimeRepo)
	mockSR.On("FindByID", ctx, showtimeID).Return(showtime, nil)

	mockLR := new(mockLockRepo)
	for _, sid := range seatIDs {
		key := "lock:seat:" + showtimeID + ":" + sid
		mockLR.On("Acquire", ctx, key, userID, 5*time.Minute).Return(true, nil)
	}

	var captured *domain.Booking
	mockBR := new(mockBookingRepo)
	mockBR.On("Create", ctx, mock.AnythingOfType("*domain.Booking")).Run(func(args mock.Arguments) {
		captured = args.Get(1).(*domain.Booking)
	}).Return(nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, seatIDs, domain.SeatStatusLocked).Return(nil)

	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Return(nil)

	mockAR := new(mockAuditRepo)

	uc := NewBookingUsecase(mockBR, mockSeatR, mockSR, new(mockMovieRepo), mockLR, mockEP, mockAR)

	booking, err := uc.LockSeat(ctx, userID, showtimeID, seatIDs)
	require.NoError(t, err)
	require.NotNil(t, booking)

	assert.NotNil(t, booking.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), *booking.ExpiresAt, time.Second)
	assert.WithinDuration(t, time.Now(), booking.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), booking.UpdatedAt, time.Second)

	assert.NotNil(t, captured)
	assert.Equal(t, *booking.ExpiresAt, *captured.ExpiresAt)

	mockSR.AssertExpectations(t)
	mockLR.AssertExpectations(t)
	mockBR.AssertExpectations(t)
	mockSeatR.AssertExpectations(t)
	mockEP.AssertExpectations(t)
}

func TestLockSeat_AcquiresRedisLockWith5mTTL(t *testing.T) {
	ctx := context.Background()
	showtime := &domain.Showtime{ID: "st-1", Price: 10}
	seatIDs := []string{"seat-1", "seat-2"}

	mockSR := new(mockShowtimeRepo)
	mockSR.On("FindByID", ctx, "st-1").Return(showtime, nil)

	mockLR := new(mockLockRepo)
	acquireCalls := 0
	for _, sid := range seatIDs {
		key := "lock:seat:st-1:" + sid
		mockLR.On("Acquire", ctx, key, "user-1", 5*time.Minute).Return(true, nil).Run(func(args mock.Arguments) {
			acquireCalls++
			ttl := args.Get(3).(time.Duration)
			assert.Equal(t, 5*time.Minute, ttl)
		})
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("Create", ctx, mock.AnythingOfType("*domain.Booking")).Return(nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, seatIDs, domain.SeatStatusLocked).Return(nil)

	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Return(nil)

	uc := NewBookingUsecase(mockBR, mockSeatR, mockSR, new(mockMovieRepo), mockLR, mockEP, new(mockAuditRepo))

	_, err := uc.LockSeat(ctx, "user-1", "st-1", seatIDs)
	require.NoError(t, err)
	assert.Equal(t, 2, acquireCalls)
}

func TestLockSeat_RollbackOnFailedLock(t *testing.T) {
	ctx := context.Background()
	showtime := &domain.Showtime{ID: "st-1", Price: 10}

	mockSR := new(mockShowtimeRepo)
	mockSR.On("FindByID", ctx, "st-1").Return(showtime, nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Acquire", ctx, "lock:seat:st-1:seat-1", "user-1", 5*time.Minute).Return(true, nil)
	mockLR.On("Acquire", ctx, "lock:seat:st-1:seat-2", "user-1", 5*time.Minute).Return(false, nil)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-1", "user-1").Return(true, nil)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-2", "user-1").Return(false, nil)

	uc := NewBookingUsecase(
		new(mockBookingRepo), new(mockSeatRepo), mockSR, new(mockMovieRepo),
		mockLR, new(mockEventPub), new(mockAuditRepo),
	)

	_, err := uc.LockSeat(ctx, "user-1", "st-1", []string{"seat-1", "seat-2"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already locked or booked")

	mockLR.AssertExpectations(t)
}

func TestExtendLock_ResetsExpiresAt(t *testing.T) {
	ctx := context.Background()
	oldExpiry := time.Now().Add(-1 * time.Minute)
	booking := &domain.Booking{
		ID:        "booking-1",
		UserID:    "user-1",
		ShowtimeID: "st-1",
		SeatIDs:   []string{"seat-1"},
		Status:    domain.BookingStatusLocked,
		ExpiresAt: &oldExpiry,
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("FindByID", ctx, "booking-1").Return(booking, nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Extend", ctx, "lock:seat:st-1:seat-1", "user-1", 5*time.Minute).Return(true, nil)

	var updatedBooking *domain.Booking
	mockBR.On("Update", ctx, mock.AnythingOfType("*domain.Booking")).Run(func(args mock.Arguments) {
		updatedBooking = args.Get(1).(*domain.Booking)
	}).Return(nil)

	uc := NewBookingUsecase(
		mockBR, new(mockSeatRepo), new(mockShowtimeRepo), new(mockMovieRepo),
		mockLR, new(mockEventPub), new(mockAuditRepo),
	)

	err := uc.ExtendLock(ctx, "user-1", "booking-1")
	require.NoError(t, err)

	require.NotNil(t, updatedBooking)
	assert.NotNil(t, updatedBooking.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), *updatedBooking.ExpiresAt, time.Second)
	assert.WithinDuration(t, time.Now(), updatedBooking.UpdatedAt, time.Second)

	mockBR.AssertExpectations(t)
	mockLR.AssertExpectations(t)
}

func TestPayBooking_ClearsExpiresAt(t *testing.T) {
	ctx := context.Background()
	futureExpiry := time.Now().Add(4 * time.Minute)
	booking := &domain.Booking{
		ID:        "booking-1",
		UserID:    "user-1",
		ShowtimeID: "st-1",
		SeatIDs:   []string{"seat-1"},
		Status:    domain.BookingStatusLocked,
		ExpiresAt: &futureExpiry,
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("FindByID", ctx, "booking-1").Return(booking, nil)

	var updatedBooking *domain.Booking
	mockBR.On("Update", ctx, mock.AnythingOfType("*domain.Booking")).Run(func(args mock.Arguments) {
		updatedBooking = args.Get(1).(*domain.Booking)
	}).Return(nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, []string{"seat-1"}, domain.SeatStatusBooked).Return(nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-1", "user-1").Return(true, nil)

	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Return(nil).Times(2)

	mockAR := new(mockAuditRepo)
	mockAR.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	uc := NewBookingUsecase(mockBR, mockSeatR, new(mockShowtimeRepo), new(mockMovieRepo), mockLR, mockEP, mockAR)

	result, err := uc.PayBooking(ctx, "user-1", "booking-1")
	require.NoError(t, err)

	require.NotNil(t, updatedBooking)
	assert.Nil(t, updatedBooking.ExpiresAt)
	assert.Equal(t, domain.BookingStatusBooked, updatedBooking.Status)
	assert.Nil(t, result.ExpiresAt)

	mockBR.AssertExpectations(t)
	mockSeatR.AssertExpectations(t)
	mockLR.AssertExpectations(t)
	mockEP.AssertExpectations(t)
	mockAR.AssertExpectations(t)
}

func TestCancelBooking_ReleasesSeatsAndLocks(t *testing.T) {
	ctx := context.Background()
	booking := &domain.Booking{
		ID:        "booking-1",
		UserID:    "user-1",
		ShowtimeID: "st-1",
		SeatIDs:   []string{"seat-1", "seat-2"},
		Status:    domain.BookingStatusLocked,
		ExpiresAt: ptr(time.Now().Add(4 * time.Minute)),
	}

	mockBR := new(mockBookingRepo)
	mockBR.On("FindByID", ctx, "booking-1").Return(booking, nil)

	var updatedBooking *domain.Booking
	mockBR.On("Update", ctx, mock.AnythingOfType("*domain.Booking")).Run(func(args mock.Arguments) {
		updatedBooking = args.Get(1).(*domain.Booking)
	}).Return(nil)

	mockSeatR := new(mockSeatRepo)
	mockSeatR.On("BulkUpdateStatus", ctx, []string{"seat-1", "seat-2"}, domain.SeatStatusAvailable).Return(nil)

	mockLR := new(mockLockRepo)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-1", "user-1").Return(true, nil)
	mockLR.On("Release", ctx, "lock:seat:st-1:seat-2", "user-1").Return(true, nil)

	mockEP := new(mockEventPub)
	mockEP.On("Publish", ctx, mock.AnythingOfType("domain.Event")).Return(nil)

	mockAR := new(mockAuditRepo)
	mockAR.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	uc := NewBookingUsecase(mockBR, mockSeatR, new(mockShowtimeRepo), new(mockMovieRepo), mockLR, mockEP, mockAR)

	err := uc.CancelBooking(ctx, "user-1", "booking-1")
	require.NoError(t, err)

	require.NotNil(t, updatedBooking)
	assert.Equal(t, domain.BookingStatusCancelled, updatedBooking.Status)
	assert.Nil(t, updatedBooking.ExpiresAt)

	mockBR.AssertExpectations(t)
	mockSeatR.AssertExpectations(t)
	mockLR.AssertExpectations(t)
	mockEP.AssertExpectations(t)
	mockAR.AssertExpectations(t)
}

func ptr(t time.Time) *time.Time {
	return &t
}
