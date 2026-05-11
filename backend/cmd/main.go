package main

import (
	"community-backend/config"
	"community-backend/internal/db"
	"community-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()
	db.Init(cfg)

	r := gin.Default()

	router.RegisterUserRoutes(r)

	r.Run(":18080")
}
