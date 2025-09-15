package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool) {
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)

	router.GET("/schedule", scheduleHandler.GetSchedules)
}
