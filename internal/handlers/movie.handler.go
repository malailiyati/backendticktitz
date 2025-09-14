package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type MovieHandler struct {
	repo *repositories.MovieRepository
}

// constructor
func NewMovieHandler(repo *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

// GetUpcomingMovies godoc
// @Summary      Get Upcoming Movies
// @Description  Get list of upcoming movies (releaseDate > today)
// @Tags         movies
// @Produce      json
// @Success      200 {array} models.MovieResponse
// @Router       /movie/upcoming [get]
func (h *MovieHandler) GetUpcomingMovies(c *gin.Context) {
	movies, err := h.repo.GetUpcomingMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

// GetPopularMovies godoc
// @Summary Get Popular Movies
// @Description Get list of popular movies ordered by popularity
// @Tags movies
// @Produce json
// @Param limit query int false "Limit number of movies (default 10)"
// @Success 200 {array} models.MovieResponse
// @Router /movie/popular [get]
func (h *MovieHandler) GetPopularMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	movies, err := h.repo.GetPopularMovies(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

// GetMoviesByFilter godoc
// @Summary Get Movies by Filter
// @Description Get movies by title and/or genre with pagination
// @Tags movies
// @Produce json
// @Param title query string false "Filter by title"
// @Param genre query string false "Filter by genre"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {array} models.MovieFilter
// @Router /movie/ [get]
func (h *MovieHandler) GetMoviesByFilter(c *gin.Context) {
	title := c.Query("title")
	genre := c.Query("genre")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	movies, err := h.repo.GetMoviesByFilter(c.Request.Context(), title, genre, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"page":    page,
		"limit":   limit,
		"data":    movies,
	})
}

// GetMovieDetail godoc
// @Summary      Get movie detail
// @Description  Get movie detail by ID
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        movie_id   path      int  true  "Movie ID"
// @Success 200 {object} models.MovieResponse
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /movies/{movie_id} [get]
func (h *MovieHandler) GetMovieDetail(c *gin.Context) {
	movieID := c.Param("movie_id")
	id, err := strconv.Atoi(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := h.repo.GetMovieDetailByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// langsung return struct MovieDetail saja
	c.JSON(http.StatusOK, movie)
}
