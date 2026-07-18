package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"github.com/naputt/cinema-ticket-backend/internal/usecase"
)

type AdminHandler struct {
	adminUsecase *usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{adminUsecase: adminUsecase}
}

// GetBookings returns all bookings with optional filters.
// @Summary List all bookings
// @Tags admin
// @Produce json
// @Param movie_id query string false "Filter by movie ID"
// @Param date query string false "Filter by date"
// @Param user_id query string false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Success 200 {array} domain.Booking
// @Router /api/v1/admin/bookings [get]
func (h *AdminHandler) GetBookings(c *gin.Context) {
	var filter usecase.BookingFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookings, err := h.adminUsecase.GetBookings(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GetAuditLogs returns audit logs with optional event type filter.
// @Summary Get audit logs
// @Tags admin
// @Produce json
// @Param event_type query string false "Filter by event type"
// @Success 200 {array} domain.AuditLog
// @Router /api/v1/admin/audit-logs [get]
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	eventType := c.Query("event_type")
	logs, err := h.adminUsecase.GetAuditLogs(c.Request.Context(), eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *AdminHandler) GetMovies(c *gin.Context) {
	movies, err := h.adminUsecase.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *AdminHandler) CreateMovie(c *gin.Context) {
	var movie domain.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.adminUsecase.CreateMovie(c.Request.Context(), &movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

func (h *AdminHandler) UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie domain.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie.ID = id
	if err := h.adminUsecase.UpdateMovie(c.Request.Context(), &movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *AdminHandler) DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminUsecase.DeleteMovie(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
