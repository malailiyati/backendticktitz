package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger(ctx *gin.Context) {
	log.Println("Start")
	start := time.Now()
	ctx.Next() // Next digunakan untuk lanjut ke middleware/handler berikutnya
	duration := time.Since(start)
	log.Printf("Durasi Request: %dus\n", duration.Microseconds())
}
