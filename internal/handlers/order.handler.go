package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/repositories"
)

type OrderHandler struct {
	repo *repositories.OrderRepository
}

func NewOrderHandler(repo *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// CreateOrder godoc
// @Summary Create new order
// @Description Create an order with schedule and seats
// @Tags orders
// @Accept json
// @Produce json
// @Param request body models.CreateOrderRequest true "Order request"
// @Success 200 {object} models.Order
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	order, err := h.repo.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})
}
