package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/pkg"
)

type MovieAdminHandler struct {
	repo *repositories.MovieAdminRepository
}

func NewMovieAdminHandler(repo *repositories.MovieAdminRepository) *MovieAdminHandler {
	return &MovieAdminHandler{repo: repo}
}

// @Summary Get all movies (Admin)
// @Description Get all movies, only accessible for admin
// @Tags Admin
// @Produce json
// @Success 200 {array} models.MovieAdmin
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Security JWTtoken
// @Router /admin/movies [get]
func (h *MovieAdminHandler) GetAllMovies(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	userClaims := claims.(pkg.Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only admin can access this endpoint",
		})
		return
	}

	movies, err := h.repo.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

// @Summary Delete a movie (Admin)
// @Description Delete movie by ID, only accessible for admin
// @Tags Admin
// @Param id path int true "Movie ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWTtoken
// @Router /admin/movies/{id} [delete]
func (h *MovieAdminHandler) DeleteMovie(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userClaims := claims.(pkg.Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only admin can access this endpoint"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid movie ID"})
		return
	}

	if err := h.repo.DeleteMovie(c.Request.Context(), id); err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie deleted successfully"})
}

func saveFile(c *gin.Context, file *multipart.FileHeader, folder, prefix string, id int) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".jfif": true}
	if !allowed[ext] {
		return "", "", fmt.Errorf("invalid file type")
	}
	if file.Size > 5<<20 {
		return "", "", fmt.Errorf("file too large")
	}

	os.MkdirAll("public/"+folder, os.ModePerm)
	newName := fmt.Sprintf("%s_%d_%d%s", prefix, id, time.Now().UnixNano(), ext)
	savePath := filepath.Join("public", folder, newName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", "", err
	}

	// return: relative path (buat DB), full path (buat rollback)
	return "/" + folder + "/" + newName, savePath, nil
}

// @Summary Patch movie (Admin)
// @Description Edit movie data (partial update, with optional poster upload)
// @Tags Admin
// @Param id path int true "Movie ID"
// @Accept multipart/form-data
// @Produce json
// @Param title formData string false "Title"
// @Param director_id formData int false "Director ID"
// @Param poster formData file false "Poster file"
// @Param background_poster formData file false "Background Poster file"
// @Param release_date formData string false "Release Date (YYYY-MM-DD)"
// @Param duration formData string false "Duration (e.g. 02:28)"
// @Param synopsis formData string false "Synopsis"
// @Param popularity formData int false "Popularity"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWTtoken
// @Router /admin/movies/{id} [patch]
// helper simpan file
func (h *MovieAdminHandler) UpdateMovie(c *gin.Context) {
	// --- Authorization check ---
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}
	userClaims := claims.(pkg.Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only admin can access"})
		return
	}

	// --- Parse movie ID ---
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid movie ID"})
		return
	}

	// --- Parse form body ---
	var form models.UpdateMovieAdminBody
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	rollbackFiles := []string{}

	// --- Normal fields ---
	if form.Title != "" {
		updates["title"] = form.Title
	}
	if form.DirectorID != "" {
		if val, err := strconv.Atoi(form.DirectorID); err == nil {
			updates["director_id"] = val
		}
	}
	if form.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", form.ReleaseDate); err == nil {
			updates["release_date"] = t
		}
	}
	if form.Duration != "" {
		updates["duration"] = form.Duration
	}
	if form.Synopsis != "" {
		updates["synopsis"] = form.Synopsis
	}
	if form.Popularity != "" {
		if val, err := strconv.Atoi(form.Popularity); err == nil {
			updates["popularity"] = val
		}
	}

	// --- Upload poster ---
	if form.Poster != nil {
		path, fullPath, err := saveFile(c, form.Poster, "posters", "poster", id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		updates["poster"] = path
		rollbackFiles = append(rollbackFiles, fullPath)
	}

	// --- Upload background ---
	if form.BackgroundPoster != nil {
		path, fullPath, err := saveFile(c, form.BackgroundPoster, "backgrounds", "bg", id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		updates["background_poster"] = path
		rollbackFiles = append(rollbackFiles, fullPath)
	}

	// --- Update DB ---
	movie, err := h.repo.UpdateMovie(c.Request.Context(), id, updates)
	if err != nil {
		// rollback file kalau DB gagal
		for _, f := range rollbackFiles {
			_ = os.Remove(f)
		}
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Movie updated successfully",
		"data":    movie,
	})
}
