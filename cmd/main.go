package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/malailiyati/backend/internal/configs"
	"github.com/malailiyati/backend/internal/routers"
)

// @title 			tickitz
// @version 		1.0
// @description 	ticketing
// @host			localhost:8080
// @basePath		/
// @securityDefinitions.apikey 	JWTtoken
// @in header
// @name Authorization
func main() {
	log.Println(os.Getenv("DBUSER")) // ini langsung kebaca dari .env

	db, err := configs.InitDB()
	if err != nil {
		log.Println("Failed to connect DB:", err)
		return
	}
	defer db.Close()

	if err := configs.TestDB(db); err != nil {
		log.Println("Ping DB failed:", err)
		return
	}
	log.Println("DB Connected")

	r := routers.InitRouter(db)

	r.Run(":8080")
}
