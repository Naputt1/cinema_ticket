package domain

import "time"

type AuditEventType string

const (
	AuditEventBookingSuccess AuditEventType = "BOOKING_SUCCESS"
	AuditEventBookingTimeout AuditEventType = "BOOKING_TIMEOUT"
	AuditEventSeatReleased   AuditEventType = "SEAT_RELEASED"
	AuditEventSystemError    AuditEventType = "SYSTEM_ERROR"
)

type AuditLog struct {
	ID        string         `json:"id" bson:"_id"`
	EventType AuditEventType `json:"event_type" bson:"event_type"`
	UserID    string         `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Detail    string         `json:"detail" bson:"detail"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
}
