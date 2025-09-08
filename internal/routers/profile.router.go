package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/malailiyati/backend/internal/handlers"
	"github.com/malailiyati/backend/internal/repositories"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	router.GET("/profile", profileHandler.GetProfile)
	router.PATCH("/profile", profileHandler.UpdateProfile)
}
