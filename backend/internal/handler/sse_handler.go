package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

type SSEHandler struct {
	mu          sync.RWMutex
	seatClients map[string]map[chan domain.Event]struct{}
	notifClients map[chan domain.Event]struct{}
}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		seatClients:  make(map[string]map[chan domain.Event]struct{}),
		notifClients: make(map[chan domain.Event]struct{}),
	}
}

type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// SeatUpdates streams seat status changes via SSE.
// @Summary Stream seat updates
// @Tags events
// @Produce text/event-stream
// @Param showtime_id query string true "Showtime ID"
// @Router /api/v1/events/seats [get]
func (h *SSEHandler) SeatUpdates(c *gin.Context) {
	showtimeID := c.Query("showtime_id")
	if showtimeID == "" {
		c.JSON(400, gin.H{"error": "showtime_id is required"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan domain.Event, 10)

	h.mu.Lock()
	if h.seatClients[showtimeID] == nil {
		h.seatClients[showtimeID] = make(map[chan domain.Event]struct{})
	}
	h.seatClients[showtimeID][ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.seatClients[showtimeID], ch)
		if len(h.seatClients[showtimeID]) == 0 {
			delete(h.seatClients, showtimeID)
		}
		h.mu.Unlock()
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-ch:
			if !ok {
				return false
			}
			data, _ := json.Marshal(SSEEvent{
				Event: string(event.Type),
				Data:  event.Payload,
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// Notifications streams notifications via SSE.
// @Summary Stream notifications
// @Tags events
// @Produce text/event-stream
// @Router /api/v1/events/notifications [get]
func (h *SSEHandler) Notifications(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan domain.Event, 10)

	h.mu.Lock()
	h.notifClients[ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.notifClients, ch)
		h.mu.Unlock()
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-ch:
			if !ok {
				return false
			}
			data, _ := json.Marshal(SSEEvent{
				Event: string(event.Type),
				Data:  event.Payload,
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func (h *SSEHandler) BroadcastSeatEvent(event domain.Event) {
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}
	showtimeID, ok := payload["showtime_id"].(string)
	if !ok {
		return
	}

	h.mu.RLock()
	clients, ok := h.seatClients[showtimeID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	h.mu.RLock()
	for ch := range clients {
		select {
		case ch <- event:
		default:
			log.Printf("SSE client buffer full, dropping event")
		}
	}
	h.mu.RUnlock()
}

func (h *SSEHandler) BroadcastNotification(event domain.Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.notifClients {
		select {
		case ch <- event:
		default:
			log.Printf("SSE notification client buffer full, dropping event")
		}
	}
}

func (h *SSEHandler) HandleEvent(event domain.Event) {
	switch {
	case event.Type == domain.EventSeatLocked || event.Type == domain.EventSeatBooked || event.Type == domain.EventSeatReleased:
		h.BroadcastSeatEvent(event)
	case event.Type == domain.EventBookingSuccess || event.Type == domain.EventBookingTimeout:
		h.BroadcastNotification(event)
	}
}
