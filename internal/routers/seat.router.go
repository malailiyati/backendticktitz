package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRepo := repositories.NewSeatRepository(db)
	seatHandler := handlers.NewSeatHandler(seatRepo)

	router.GET("/seats", seatHandler.GetSoldSeats)
}
