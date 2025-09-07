package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type MovieFilterHandler struct {
	repo *repositories.MovieFilterRepository
}

func NewMovieFilterHandler(repo *repositories.MovieFilterRepository) *MovieFilterHandler {
	return &MovieFilterHandler{repo: repo}
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
// @Router /movie/filter [get]
func (h *MovieFilterHandler) GetMoviesByFilter(c *gin.Context) {
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
