package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	router.POST("/orders", orderHandler.CreateOrder)
}
