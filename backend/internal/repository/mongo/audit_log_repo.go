package mongo

import (
	"context"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type auditLogRepository struct {
	coll *mongo.Collection
}

func NewAuditLogRepository(db *mongo.Database) domain.AuditLogRepository {
	return &auditLogRepository{
		coll: db.Collection("audit_logs"),
	}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	_, err := r.coll.InsertOne(ctx, log)
	return err
}

func (r *auditLogRepository) BulkCreate(ctx context.Context, logs []domain.AuditLog) error {
	docs := make([]interface{}, len(logs))
	for i := range logs {
		docs[i] = logs[i]
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *auditLogRepository) FindAll(ctx context.Context, filter map[string]interface{}) ([]domain.AuditLog, error) {
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
