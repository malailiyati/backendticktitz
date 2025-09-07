package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	origin := ctx.GetHeader("Origin")

	if origin != "" {
		// izinkan origin mana pun yang minta
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
