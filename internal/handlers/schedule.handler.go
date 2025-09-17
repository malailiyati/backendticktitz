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

// GetSchedulesByFilter godoc
// @Summary      Get schedules by filter
// @Description  Get schedules based on date, time, location, and movie
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        movie_id   path  int    true "Movie ID"
// @Param        date        query   string false "Schedule date (YYYY-MM-DD)"
// @Param        time_id     query   int    false "Time ID"
// @Param        location_id query   int    false "Location ID"
// @Success      200 {object} map[string]interface{} "success response with schedules"
// @Failure      400 {object} map[string]interface{} "invalid input"
// @Failure      500 {object} map[string]interface{} "internal server error"
// @Router       /movie/{movie_id}/schedule [get]
func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	movieID, _ := strconv.Atoi(c.Param("movie_id"))
	date := c.Query("date")
	timeID, _ := strconv.Atoi(c.Query("time_id"))
	locationID, _ := strconv.Atoi(c.Query("location_id"))

	schedules, err := h.repo.GetSchedules(c.Request.Context(), date, timeID, locationID, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}
