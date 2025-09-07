package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/internal/utils"
)

type MoviePopularHandler struct {
	repo *repositories.MoviePopularRepository
}

func NewMoviePopularHandler(repo *repositories.MoviePopularRepository) *MoviePopularHandler {
	return &MoviePopularHandler{repo: repo}
}

// GetPopularMovies godoc
// @Summary Get Popular Movies
// @Description Get list of popular movies ordered by popularity
// @Tags movies
// @Produce json
// @Param limit query int false "Limit number of movies (default 10)"
// @Success 200 {array} MovieResponse
// @Router /movie/popular [get]
func (h *MoviePopularHandler) GetPopularMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	movies, err := h.repo.GetPopularMovies(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	var response []MovieResponse
	for _, m := range movies {
		response = append(response, MovieResponse{
			ID:               m.ID,
			Title:            m.Title,
			DirectorID:       m.DirectorID,
			Poster:           m.Poster,
			BackgroundPoster: m.BackgroundPoster,
			ReleaseDate:      m.ReleaseDate,
			Duration:         utils.FormatIntervalToText(m.Duration),
			Synopsis:         m.Synopsis,
			Popularity:       m.Popularity,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
