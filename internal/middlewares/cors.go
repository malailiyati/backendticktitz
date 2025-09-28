package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		whitelist := []string{
			"http://127.0.0.1:5500",
			"http://127.0.0.1:3001",
			"http://localhost:5173",
			"http://localhost", // nginx prod
		}
		origin := ctx.GetHeader("Origin")
		if slices.Contains(whitelist, origin) {
			ctx.Header("Access-Control-Allow-Origin", origin)
		}
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// package middlewares

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func CORSMiddleware(ctx *gin.Context) {
// 	origin := ctx.GetHeader("Origin")

// 	if origin != "" {
// 		// izinkan origin mana pun yang minta
// 		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 	} else {
// 		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 	}

// 	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
// 	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

// 	if ctx.Request.Method == http.MethodOptions {
// 		ctx.AbortWithStatus(http.StatusNoContent)
// 		return
// 	}

// 	ctx.Next()
// }
