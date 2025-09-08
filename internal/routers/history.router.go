package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitHistoryRouter(router *gin.Engine, db *pgxpool.Pool) {
	historyRepo := repositories.NewHistoryRepository(db)
	historyHandler := handlers.NewHistoryHandler(historyRepo)

	router.GET("/history", historyHandler.GetHistory)
}
