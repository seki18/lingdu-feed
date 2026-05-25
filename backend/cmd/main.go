package main

import (
	"github.com/seki18/lingdu-feed/config"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/router"
	"github.com/seki18/lingdu-feed/internal/scheduler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()
	common.Init(cfg)
	common.InitRedis(cfg)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.RegisterRoutes(r)

	// Start background score scheduler: full update on startup, then every 1 minute
	go scheduler.RunScoreScheduler()

	r.Run(":18080")
}
