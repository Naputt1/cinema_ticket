package usecase

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminUsecase struct {
	bookingRepo domain.BookingRepository
	auditRepo   domain.AuditLogRepository
}

func NewAdminUsecase(
	bookingRepo domain.BookingRepository,
	auditRepo domain.AuditLogRepository,
) *AdminUsecase {
	return &AdminUsecase{
		bookingRepo: bookingRepo,
		auditRepo:   auditRepo,
	}
}

type BookingFilter struct {
	MovieID string `form:"movie_id"`
	Date    string `form:"date"`
	UserID  string `form:"user_id"`
	Status  string `form:"status"`
}

func (uc *AdminUsecase) GetBookings(ctx context.Context, filter BookingFilter) ([]domain.Booking, error) {
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
	return uc.bookingRepo.FindAll(ctx, f)
}

func (uc *AdminUsecase) GetAuditLogs(ctx context.Context, eventType string) ([]domain.AuditLog, error) {
	f := bson.M{}
	if eventType != "" {
		f["event_type"] = eventType
	}
	return uc.auditRepo.FindAll(ctx, f)
}
