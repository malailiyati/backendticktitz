package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/malailiyati/backend/pkg"
)

func VerifyToken(ctx *gin.Context) {
	// ambil token dari header
	bearerToken := ctx.GetHeader("Authorization")
	// Bearer token
	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Silahkan login terlebih dahulu",
		})
		return
	}
	token := parts[1]

	var claims pkg.Claims
	if err := claims.VerifyToken(token); err != nil {
		if strings.Contains(err.Error(), jwt.ErrTokenInvalidIssuer.Error()) {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Silahkan login kembali",
			})
			return
		}
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Silahkan login kembali",
			})
			return
		}
		fmt.Println(jwt.ErrTokenExpired)
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}
