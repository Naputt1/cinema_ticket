package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/naputt/cinema-ticket-backend/internal/usecase"
)

type BookingHandler struct {
	bookingUsecase *usecase.BookingUsecase
}

func NewBookingHandler(bookingUsecase *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUsecase: bookingUsecase}
}

// LockSeat locks seats for booking.
// @Summary Lock seats
// @Tags booking
// @Accept json
// @Produce json
// @Param request body object{showtime_id=string,seat_ids=[]string} true "Lock request"
// @Success 200 {object} domain.Booking
// @Failure 409 {object} map[string]string
// @Router /api/v1/bookings/lock [post]
func (h *BookingHandler) LockSeat(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		ShowtimeID string   `json:"showtime_id" binding:"required"`
		SeatIDs    []string `json:"seat_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.bookingUsecase.LockSeat(c.Request.Context(), userID, req.ShowtimeID, req.SeatIDs)
	if err != nil {
		if strings.Contains(err.Error(), "already locked") || strings.Contains(err.Error(), "already booked") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// PayBooking simulates payment and confirms booking.
// @Summary Confirm payment
// @Tags booking
// @Accept json
// @Produce json
// @Param request body object{booking_id=string} true "Payment request"
// @Success 200 {object} domain.Booking
// @Router /api/v1/bookings/pay [post]
func (h *BookingHandler) PayBooking(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		BookingID string `json:"booking_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing booking_id"})
		return
	}

	booking, err := h.bookingUsecase.PayBooking(c.Request.Context(), userID, req.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// CancelMyBooking cancels a booking by ID (for LOCKED bookings only).
func (h *BookingHandler) CancelMyBooking(c *gin.Context) {
	userID := c.GetString("user_id")
	bookingID := c.Param("id")

	if err := h.bookingUsecase.CancelBooking(c.Request.Context(), userID, bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// CancelBooking cancels a locked booking.
// @Summary Cancel booking
// @Tags booking
// @Accept json
// @Produce json
// @Param request body object{booking_id=string} true "Cancel request"
// @Success 200 {object} map[string]string
// @Router /api/v1/bookings/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		BookingID string `json:"booking_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing booking_id"})
		return
	}

	if err := h.bookingUsecase.CancelBooking(c.Request.Context(), userID, req.BookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// GetMyBookings returns current user's bookings.
// @Summary Get user bookings
// @Tags booking
// @Produce json
// @Success 200 {array} domain.Booking
// @Router /api/v1/bookings/my [get]
func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetString("user_id")
	bookings, err := h.bookingUsecase.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GetBooking returns booking details.
// @Summary Get booking details
// @Tags booking
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} domain.Booking
// @Router /api/v1/bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	id := c.Param("id")
	booking, err := h.bookingUsecase.GetBooking(c.Request.Context(), id)
	if err != nil || booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	c.JSON(http.StatusOK, booking)
}
