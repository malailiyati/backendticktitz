package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/middlewares"
	"github.com/malailiyati/backend/internal/models"

	docs "github.com/malailiyati/backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.MyLogger)
	router.Use(middlewares.CORSMiddleware)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// config := cors.Config{
	// 	AllowOrigins: []string{"http://127.0.0.1:5500", "http://127.0.0.1:3001"},
	// 	AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders: []string{"Authorization", "Content-Type"},
	// }

	// router.Use(cors.New(config))

	router.Static("/img", "public")

	// setup routing
	// InitPingRouter(router)
	InitAuthRouter(router, db)
	InitMovieRouter(router, db)
	InitScheduleRouter(router, db)
	InitSeatRouter(router, db)
	InitOrderRouter(router, db)
	InitProfileRouter(router, db)
	InitHistoryRouter(router, db)
	InitMovieAdminRouter(router, db)

	// catch all route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			Message: "Rute Salah",
			Status:  "Rute Tidak Ditemukan",
		})
	})

	return router
}
