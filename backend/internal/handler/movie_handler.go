package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naputt/cinema-ticket-backend/internal/usecase"
)

type MovieHandler struct {
	movieUsecase *usecase.MovieUsecase
}

func NewMovieHandler(movieUsecase *usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{movieUsecase: movieUsecase}
}

// GetMovies lists all movies.
// @Summary List all movies
// @Tags movies
// @Produce json
// @Success 200 {array} domain.Movie
// @Router /api/v1/movies [get]
func (h *MovieHandler) GetMovies(c *gin.Context) {
	movies, err := h.movieUsecase.GetMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetMovie returns a movie by ID.
// @Summary Get movie details
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} domain.Movie
// @Router /api/v1/movies/{id} [get]
func (h *MovieHandler) GetMovie(c *gin.Context) {
	id := c.Param("id")
	movie, err := h.movieUsecase.GetMovie(c.Request.Context(), id)
	if err != nil || movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// GetShowtimes returns showtimes for a movie.
// @Summary Get showtimes by movie
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {array} domain.Showtime
// @Router /api/v1/movies/{id}/showtimes [get]
func (h *MovieHandler) GetShowtimes(c *gin.Context) {
	id := c.Param("id")
	showtimes, err := h.movieUsecase.GetShowtimes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)
}

// GetSeats returns seat map for a showtime.
// @Summary Get seat map
// @Tags showtimes
// @Produce json
// @Param id path string true "Showtime ID"
// @Success 200 {array} domain.Seat
// @Router /api/v1/showtimes/{id}/seats [get]
func (h *MovieHandler) GetSeats(c *gin.Context) {
	id := c.Param("id")
	seats, err := h.movieUsecase.GetSeats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}
