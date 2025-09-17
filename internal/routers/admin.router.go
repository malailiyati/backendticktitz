package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/middlewares"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitMovieAdminRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieAdminRepo := repositories.NewMovieAdminRepository(db)
	movieAdminHandler := handlers.NewMovieAdminHandler(movieAdminRepo)

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(rdb, "admin"))
	{
		admin.GET("/movies", movieAdminHandler.GetAllMovies)
		admin.DELETE("/movies/:id", movieAdminHandler.DeleteMovie)
		admin.PATCH("/movies/:id", movieAdminHandler.UpdateMovie)
		admin.POST("/movies", movieAdminHandler.CreateMovie)
	}

}
