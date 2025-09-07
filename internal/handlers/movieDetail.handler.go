package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type MovieDetailHandler struct {
	repo *repositories.MovieDetailRepository
}

func NewMovieDetailHandler(repo *repositories.MovieDetailRepository) *MovieDetailHandler {
	return &MovieDetailHandler{repo: repo}
}

// GetMovieDetail godoc
// @Summary Get movie detail
// @Description Get detailed information for a movie (genres, casts, director)
// @Tags movies
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Success 200 {object} models.MovieDetail
// @Router /movie/detail [get]
func (h *MovieDetailHandler) GetMovieDetail(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Query("movie_id"))
	if err != nil || movieID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid movie_id"})
		return
	}

	detail, err := h.repo.GetMovieDetail(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": detail})
}
