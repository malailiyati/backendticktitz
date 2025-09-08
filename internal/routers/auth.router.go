package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("login", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/logout", authHandler.Logout)
}
