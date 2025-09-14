package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/repositories"
)

type HistoryHandler struct {
	repo *repositories.HistoryRepository
}

func NewHistoryHandler(repo *repositories.HistoryRepository) *HistoryHandler {
	return &HistoryHandler{repo: repo}
}

// GetHistory godoc
// @Summary Get order history
// @Description Get all order history for a user
// @Tags history
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} models.OrderHistory
// @Security JWTtoken
// @Router /user/history [get]
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid user_id"})
		return
	}

	history, err := h.repo.GetOrderHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": history})
}
