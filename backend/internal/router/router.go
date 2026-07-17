package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"github.com/naputt/cinema-ticket-backend/internal/handler"
	"github.com/naputt/cinema-ticket-backend/internal/middleware"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Movie   *handler.MovieHandler
	Booking *handler.BookingHandler
	Admin   *handler.AdminHandler
	SSE     *handler.SSEHandler
}

func Setup(r *gin.Engine, h *Handlers, jwtSecret string, frontendURL string) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	api.GET("/auth/google/login", h.Auth.GoogleLogin)
	api.GET("/auth/google/callback", h.Auth.GoogleCallback)
	api.POST("/auth/refresh", h.Auth.RefreshToken)
	api.GET("/auth/logout", h.Auth.Logout)

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.GET("/auth/me", h.Auth.GetMe)

		auth.GET("/movies", h.Movie.GetMovies)
		auth.GET("/movies/:id", h.Movie.GetMovie)
		auth.GET("/movies/:id/showtimes", h.Movie.GetShowtimes)
		auth.GET("/showtimes/:id/seats", h.Movie.GetSeats)

		auth.POST("/bookings/lock", h.Booking.LockSeat)
		auth.POST("/bookings/pay", h.Booking.PayBooking)
		auth.POST("/bookings/cancel", h.Booking.CancelBooking)
		auth.GET("/bookings/my", h.Booking.GetMyBookings)
		auth.GET("/bookings/:id", h.Booking.GetBooking)

		auth.GET("/events/seats", h.SSE.SeatUpdates)
		auth.GET("/events/notifications", h.SSE.Notifications)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware(domain.RoleAdmin))
	{
		admin.GET("/bookings", h.Admin.GetBookings)
		admin.GET("/audit-logs", h.Admin.GetAuditLogs)
		admin.GET("/movies", h.Admin.GetMovies)
		admin.POST("/movies", h.Admin.CreateMovie)
		admin.PUT("/movies/:id", h.Admin.UpdateMovie)
		admin.DELETE("/movies/:id", h.Admin.DeleteMovie)
	}
}
