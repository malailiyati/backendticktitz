package handlers

import (
	"net/http"

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
// @Tags profile
// @Produce json
// @Success 200 {array} models.OrderHistory
// @Security JWTtoken
// @Router /user/history [get]
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User ID tidak ditemukan di token",
		})
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "User ID tidak valid",
		})
		return
	}

	history, err := h.repo.GetOrderHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": history})
}
