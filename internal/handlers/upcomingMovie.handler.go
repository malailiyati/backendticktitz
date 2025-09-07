package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/internal/utils"
)

type MovieHandler struct {
	repo *repositories.MovieRepository
}

// constructor
func NewUpcomingMovieHandler(repo *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

// GetUpcomingMovies godoc
// @Summary      Get Upcoming Movies
// @Description  Get list of upcoming movies (releaseDate > today)
// @Tags         movies
// @Produce      json
// @Success      200 {array} MovieResponse
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

	var response []MovieResponse
	for _, m := range movies {
		response = append(response, MovieResponse{
			ID:               m.ID,
			Title:            m.Title,
			DirectorID:       m.DirectorID,
			Poster:           m.Poster,
			BackgroundPoster: m.BackgroundPoster,
			ReleaseDate:      m.ReleaseDate,
			Duration:         utils.FormatIntervalToText(m.Duration), // "2 jam 30 menit"
			Synopsis:         m.Synopsis,
			Popularity:       m.Popularity,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// response struct khusus API
type MovieResponse struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	DirectorID       int       `json:"director_id"`
	Poster           string    `json:"poster"`
	BackgroundPoster string    `json:"background_poster"`
	ReleaseDate      time.Time `json:"release_date"`
	Duration         string    `json:"duration"`
	Synopsis         string    `json:"synopsis"`
	Popularity       int       `json:"popularity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
