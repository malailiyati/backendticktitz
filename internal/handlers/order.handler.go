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
// @Security JWTtoken
// @Router /user/orders [post]
func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req models.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Inject user_id dari JWT (biar user ga bisa order atas nama orang lain)
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Silahkan login terlebih dahulu",
		})
		return
	}
	req.UserID = userID

	// Panggil repository
	order, err := h.repo.CreateOrder(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Response sukses
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    order,
	})
}
