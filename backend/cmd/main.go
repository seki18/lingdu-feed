package main

import (
"log"
"os"
"os/signal"
"syscall"
"time"

"github.com/seki18/lingdu-feed/config"
"github.com/seki18/lingdu-feed/internal/cache"
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

// Start periodic stats cache sync to DB every 30 seconds
go func() {
ticker := time.NewTicker(30 * time.Second)
defer ticker.Stop()
for range ticker.C {
cache.SyncAllToDB()
}
}()

// Graceful shutdown: flush stats cache one final time
go func() {
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
log.Println("[Main] Shutting down, flushing stats cache...")
cache.SyncAllToDB()
os.Exit(0)
}()

r.Run(":18080")
}
