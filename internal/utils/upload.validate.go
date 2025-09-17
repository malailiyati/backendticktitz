package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveFile(c *gin.Context, file *multipart.FileHeader, folder, prefix string, id int) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".jfif": true}
	if !allowed[ext] {
		return "", "", fmt.Errorf("invalid file type")
	}
	log.Println(file.Size)
	if file.Size > 500*1024 {
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
