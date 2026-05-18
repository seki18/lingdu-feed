package main

import (
	"github.com/seki18/lingdu-feed/config"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()
	common.Init(cfg)

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

	router.UserRoutes(r)
	router.FollowRoutes(r)
	router.FeedRoutes(r)
	router.PostRoutes(r)
	router.CommentRoutes(r)
	router.PraiseRoutes(r)
	router.CollectionRoutes(r)
	router.InteractionStatusRoutes(r)

	r.Run(":18080")
}
