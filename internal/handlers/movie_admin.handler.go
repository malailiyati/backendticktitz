package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/repositories"
	"github.com/malailiyati/backend/internal/utils"
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
			updates["releasedate"] = t
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
		path, fullPath, err := utils.SaveFile(c, form.Poster, "posters", "poster", id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		updates["poster"] = path
		rollbackFiles = append(rollbackFiles, fullPath)
	}

	// --- Upload background ---
	if form.BackgroundPoster != nil {
		path, fullPath, err := utils.SaveFile(c, form.BackgroundPoster, "backgrounds", "bg", id)
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

// CreateMovie godoc
// @Summary      Create new movie (admin only)
// @Description  Admin can create a new movie with poster & background upload
// @Tags         Admin
// @Accept       multipart/form-data
// @Produce      json
// @Param        title             formData string true  "Movie Title"
// @Param        synopsis          formData string true  "Synopsis"
// @Param        release_date      formData string true  "Release Date (YYYY-MM-DD)"
// @Param        duration          formData string true  "Duration (HH:MM:SS)"
// @Param        director_id       formData int    true  "Director ID"
// @Param        popularity        formData int    false "Popularity"
// @Param        poster            formData file   false "Poster file"
// @Param        background_poster formData file   false "Background poster file"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Security     JWTtoken
// @Router       /admin/movies [post]
func (h *MovieAdminHandler) CreateMovie(c *gin.Context) {
	var form models.CreateMovieAdminBody
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// --- Parse release_date ---
	releaseDate, err := time.Parse("2006-01-02", form.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid release_date format"})
		return
	}

	// --- Parse duration (format HH:MM:SS) ---
	var duration pgtype.Interval
	if form.Duration != "" {
		if d, err := time.ParseDuration(form.Duration); err == nil {
			duration = pgtype.Interval{Microseconds: d.Microseconds(), Valid: true}
		}
	}

	movie := models.Movie{
		Title:       form.Title,
		Synopsis:    form.Synopsis,
		ReleaseDate: releaseDate,
		Duration:    duration,
		DirectorID:  form.DirectorID,
		Popularity:  form.Popularity,
	}

	rollbackFiles := []string{}

	// --- Upload poster ---
	if form.Poster != nil {
		path, fullPath, err := utils.SaveFile(c, form.Poster, "posters", "poster", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		movie.Poster = path
		rollbackFiles = append(rollbackFiles, fullPath)
	}

	// --- Upload background poster ---
	if form.BackgroundPoster != nil {
		path, fullPath, err := utils.SaveFile(c, form.BackgroundPoster, "backgrounds", "bg", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		movie.BackgroundPoster = path
		rollbackFiles = append(rollbackFiles, fullPath)
	}

	// --- Insert DB ---
	newMovie, err := h.repo.CreateMovie(c.Request.Context(), movie)
	if err != nil {
		// rollback kalau DB gagal
		for _, f := range rollbackFiles {
			_ = os.Remove(f)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Movie created successfully",
		"data":    newMovie,
	})
}
