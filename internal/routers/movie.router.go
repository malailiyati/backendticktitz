package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	moviesRepo := repositories.NewMovieRepository(db, rdb)
	moviesHandler := handlers.NewMovieHandler(moviesRepo)

	movieRouter := router.Group("/movie")
	movieRouter.GET("/upcoming", moviesHandler.GetUpcomingMovies)
	movieRouter.GET("/popular", moviesHandler.GetPopularMovies)
	movieRouter.GET("", moviesHandler.GetMoviesByFilter)
	movieRouter.GET("/:movie_id", moviesHandler.GetMovieDetail)
}
