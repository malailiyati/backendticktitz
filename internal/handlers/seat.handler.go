package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type SeatHandler struct {
	repo *repositories.SeatRepository
}

func NewSeatHandler(repo *repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{repo: repo}
}

// GetAvailableSeats godoc
// @Summary Get sold seats by schedule
// @Description Get all sold seats for a schedule
// @Tags seats
// @Produce json
// @Param schedule_id query int true "Schedule ID"
// @Success 200 {array} models.Seat
// @Router /seats [get]
func (h *SeatHandler) GetSoldSeats(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Query("schedule_id"))
	if err != nil || scheduleID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid schedule_id"})
		return
	}

	seats, err := h.repo.GetAvailableSeats(c.Request.Context(), scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": seats})
}
