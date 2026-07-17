package domain

import "context"

type EventType string

const (
	EventSeatLocked   EventType = "seat.locked"
	EventSeatBooked   EventType = "seat.booked"
	EventSeatReleased EventType = "seat.released"
	EventBookingSuccess EventType = "booking.success"
	EventBookingTimeout EventType = "booking.timeout"
	EventSystemError   EventType = "system.error"
	EventNotification  EventType = "notification"
)

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type EventConsumer interface {
	Consume(ctx context.Context, handler func(Event) error) error
}
