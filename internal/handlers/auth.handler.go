package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/internal/utils"
	"github.com/malailiyati/backend/pkg"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Login godoc
// @Summary      Login user
// @Description  Login dengan email dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.Login  true  "request"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Router       /auth/login [post]
func (a *AuthHandler) Login(ctx *gin.Context) {
	var body models.Login
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Data tidak valid",
		})
		return
	}

	log.Println("Login attempt:", body.Email)

	// cari user by email
	user, err := a.ar.GetUserWithPasswordAndRole(ctx.Request.Context(), body.Email)
	if err != nil {
		log.Println("User not found / DB error:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Email atau password salah",
		})
		return
	}

	// bandingkan password
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.Password, user.Password)
	if err != nil {
		log.Println("Password hash error:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Email atau password salah",
		})
		return
	}

	if !isMatched {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Email atau password salah",
		})
		return
	}

	// generate JWT token
	claims := pkg.NewJWTClaims(user.Id, user.Role)
	token, err := claims.GenToken()
	if err != nil {
		log.Println("JWT error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

// Register godoc
// @Summary      Register user baru
// @Description  Membuat akun baru
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.Register  true  "Register request (role optional, default user)"
// @Success      201   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Router       /auth/register [post]
func (a *AuthHandler) Register(ctx *gin.Context) {
	var body models.Register
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Data tidak valid",
		})
		return
	}

	// validasi
	if err := utils.ValidateRegister(models.UserAuth{
		Email:    body.Email,
		Password: body.Password,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// hash password
	hc := pkg.NewHashConfig()
	hc.UseRecommended()
	hashedPass, err := hc.GenHash(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat password",
		})
		return
	}

	// set default role kalau kosong
	if body.Role == nil {
		defaultRole := "user"
		body.Role = &defaultRole
	}

	// simpan user
	user, err := a.ar.CreateUser(ctx.Request.Context(), body.Email, hashedPass, body.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat user",
		})
		return
	}

	// hanya balikin email
	resp := models.RegisterResponse{
		Email: user.Email,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    resp,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Menghapus session (stateless di JWT)
// @Tags auth
// @Produce json
// @Success 200 {object} models.Response
// @Security JWTtoken
// @Router /auth/logout [post]
func (a *AuthHandler) Logout(ctx *gin.Context, rdb *redis.Client) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token tidak ditemukan",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	var claims pkg.Claims
	if err := claims.VerifyToken(tokenString); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Token tidak valid",
		})
		return
	}

	exp := time.Until(claims.ExpiresAt.Time)

	// simpan token ke Redis blacklist dengan TTL = sisa expired
	if err := rdb.Set(ctx, "blacklist:"+tokenString, true, exp).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal blacklist token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil, token di-blacklist",
	})
}
