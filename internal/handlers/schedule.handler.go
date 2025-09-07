package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type ScheduleHandler struct {
	repo *repositories.ScheduleRepository
}

func NewScheduleHandler(repo *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{repo: repo}
}

// GetSchedulesByMovie godoc
// @Summary Get schedules by movie ID
// @Description Get schedule list (date, time, location, cinema, price) for a movie
// @Tags schedules
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Success 200 {array} models.ScheduleDetail
// @Router /schedule [get]
func (h *ScheduleHandler) GetSchedulesByMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Query("movie_id"))
	if err != nil || movieID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid movie_id"})
		return
	}

	schedules, err := h.repo.GetSchedulesByMovie(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": schedules})
}
