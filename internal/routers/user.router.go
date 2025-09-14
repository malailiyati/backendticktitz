package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/middlewares"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	// history
	historyRepo := repositories.NewHistoryRepository(db)
	historyHandler := handlers.NewHistoryHandler(historyRepo)

	// profile
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	// User routes (hanya bisa diakses user yang login)
	user := router.Group("/user")
	user.Use(middlewares.AuthMiddleware("user"))
	{
		user.POST("/orders", orderHandler.CreateOrder)  // buat order tiket
		user.GET("/history", historyHandler.GetHistory) // lihat riwayat order
		user.GET("/profile", profileHandler.GetProfile)
		user.PATCH("/profile", profileHandler.UpdateProfile)
		user.PUT("/password", profileHandler.UpdatePassword)
	}
}
