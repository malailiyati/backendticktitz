package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/internal/utils"
	"github.com/malailiyati/backend/pkg"
)

type ProfileHandler struct {
	repo *repositories.ProfileRepository
}

func NewProfileHandler(repo *repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{repo: repo}
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update firstname, lastname, phone, profile_picture
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Param user_id query int true "User ID"
// @Param first_name formData string false "First Name"
// @Param last_name formData string false "Last Name"
// @Param phone formData string false "Phone"
// @Param profile_picture formData file false "Profile Picture"
// @Success 200 {object} models.Profile
// @Security JWTtoken
// @Router /user/profile [patch]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid user_id"})
		return
	}

	var req models.ProfileBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["firstname"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["lastname"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	// upload file
	if req.ProfilePictureFile != nil {
		filename := fmt.Sprintf("public/profile_pictures/%d_%s", userID, req.ProfilePictureFile.Filename)
		if err := c.SaveUploadedFile(req.ProfilePictureFile, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to save file"})
			return
		}
		updates["profile_picture"] = filename
	}

	updated, err := h.repo.UpdateProfile(c.Request.Context(), userID, updates)
	if err != nil {
		// rollback kalau DB gagal tapi file sudah tersimpan
		if req.ProfilePictureFile != nil {
			_ = os.Remove(updates["profile_picture"].(string))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile information by user_id (join users + profile)
// @Tags profile
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {object} models.ProfileResponse
// @Security JWTtoken
// @Router /user/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid user_id"})
		return
	}

	profile, err := h.repo.GetProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// UpdatePassword godoc
// @Summary Ubah password user
// @Tags User
// @Security JWTtoken
// @Accept json
// @Produce json
// @Param body body models.UpdatePasswordRequest true "Password lama & baru"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string{error=string}
// @Failure 401 {object} map[string]string{error=string}
// @Failure 404 {object} map[string]string{error=string}
// @Router /user/password [put]
func (h *ProfileHandler) UpdatePassword(ctx *gin.Context) {
	var req models.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// cek confirm password
	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password baru tidak sama"})
		return
	}

	// validasi strength password baru
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ambil user id dari JWT claims
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Silahkan login terlebih dahulu"})
		return
	}

	// ambil data user dari DB
	user, err := h.repo.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	fmt.Println("DEBUG: password DB =", user.Password)
	fmt.Println("DEBUG: password input =", req.OldPassword)

	// verifikasi password lama
	hc := pkg.NewHashConfig()
	ok, err := hc.CompareHashAndPassword(req.OldPassword, user.Password)
	if err != nil || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
		return
	}

	// hash password baru
	newHash, err := hc.GenHash(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hash password baru"})
		return
	}

	// update password di DB
	if err := h.repo.UpdatePassword(ctx.Request.Context(), userID, newHash); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
