package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/pkg"
)

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token tidak valid",
			})
			return
		}

		// simpan claims ke context
		ctx.Set("claims", claims)

		// cek role kalau ada rules
		if len(allowedRoles) > 0 {
			role := claims.Role
			allowed := false
			for _, r := range allowedRoles {
				if role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"error":   "Forbidden",
				})
				return
			}
		}

		ctx.Next()
	}
}
