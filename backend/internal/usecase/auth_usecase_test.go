package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"github.com/naputt/cinema-ticket-backend/internal/middleware"
	"github.com/naputt/cinema-ticket-backend/pkg/config"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) FindAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, u *domain.User) error {
	return m.Called(ctx, u).Error(0)
}

func TestRefreshAccessToken_AccessTokenHasCorrectExpiry(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 168 * time.Hour,
		},
	}

	mockUR := new(mockUserRepo)
	mockUR.On("FindByID", mock.Anything, "user-1").Return(&domain.User{
		ID:    "user-1",
		Email: "test@test.com",
		Name:  "Test User",
		Role:  domain.RoleUser,
	}, nil)

	uc := NewAuthUsecase(mockUR, cfg)

	now := time.Now()
	refreshClaims := middleware.Claims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(168 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte("test-secret"))
	require.NoError(t, err)

	accessToken, err := uc.RefreshAccessToken(context.Background(), refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	accessClaims := &middleware.Claims{}
	parsed, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
	assert.WithinDuration(t, now.Add(15*time.Minute), accessClaims.ExpiresAt.Time, 2*time.Second)
	assert.Equal(t, "user-1", accessClaims.UserID)
	assert.Equal(t, domain.RoleUser, accessClaims.Role)
}

func TestRefreshAccessToken_ExpiredToken_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 168 * time.Hour,
		},
	}

	uc := NewAuthUsecase(new(mockUserRepo), cfg)

	expiredClaims := middleware.Claims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte("test-secret"))
	require.NoError(t, err)

	_, err = uc.RefreshAccessToken(context.Background(), expiredToken)
	assert.Error(t, err)
}

func TestGetAccessExpiry_ReturnsConfiguredDuration(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessExpiry:  20 * time.Minute,
			RefreshExpiry: 336 * time.Hour,
		},
	}

	uc := NewAuthUsecase(new(mockUserRepo), cfg)
	assert.Equal(t, 20*time.Minute, uc.GetAccessExpiry())
	assert.Equal(t, 336*time.Hour, uc.GetRefreshExpiry())
}
