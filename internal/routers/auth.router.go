package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("/login", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	// inject rdb ke handler logout
	authRouter.POST("/logout", func(c *gin.Context) {
		authHandler.Logout(c, rdb)
	})
}
