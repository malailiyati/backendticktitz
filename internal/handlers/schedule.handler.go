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
// @Description  Get schedules based on date, time, location, and cinema
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        date        query   string true "Schedule date (YYYY-MM-DD)"
// @Param        time_id     query   int    true "Time ID"
// @Param        location_id query   int    true "Location ID"
// @Param        cinema_id   query   int    true "Cinema ID"
// @Success      200 {object} map[string]interface{} "success response with schedules"
// @Failure      400 {object} map[string]interface{} "invalid input"
// @Failure      500 {object} map[string]interface{} "internal server error"
// @Router       /schedule [get]
func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	date := c.Query("date")
	timeID, _ := strconv.Atoi(c.Query("time_id"))
	locationID, _ := strconv.Atoi(c.Query("location_id"))
	cinemaID, _ := strconv.Atoi(c.Query("cinema_id"))

	if date == "" || timeID == 0 || locationID == 0 || cinemaID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "date, time_id, location_id, and cinema_id are required",
		})
		return
	}

	schedules, err := h.repo.GetSchedules(c.Request.Context(), date, timeID, locationID, cinemaID)
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
