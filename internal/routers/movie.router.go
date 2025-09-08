package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	// upcoming
	upcomingRepo := repositories.NewUpcomingMovieRepository(db)
	upcomingHandler := handlers.NewUpcomingMovieHandler(upcomingRepo)

	// popular
	popularRepo := repositories.NewMoviePopularRepository(db)
	popularHandler := handlers.NewMoviePopularHandler(popularRepo)

	// filter
	filterRepo := repositories.NewMovieFilterRepository(db)
	filterHandler := handlers.NewMovieFilterHandler(filterRepo)

	// repo & handler untuk detail
	detailRepo := repositories.NewMovieDetailRepository(db)
	detailHandler := handlers.NewMovieDetailHandler(detailRepo)

	movieRouter := router.Group("/movie")
	movieRouter.GET("/upcoming", upcomingHandler.GetUpcomingMovies)
	movieRouter.GET("/popular", popularHandler.GetPopularMovies)
	movieRouter.GET("/filter", filterHandler.GetMoviesByFilter)
	movieRouter.GET("/detail", detailHandler.GetMovieDetail)
}
