package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"bsf-engine/receiver"
	"bsf-engine/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load(".env.local")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	redisAddr := os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr + ":" + redisPort,
		PoolSize:     10,
		MinIdleConns: 5,
	})
	defer rdb.Close()

	if err := receiver.LoadToRedis(rdb, "assets/bahrain_roads.json"); err != nil {
		log.Printf("failed to preload data: %v", err)
	}

	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Reloading road data into Redis...")
			if err := receiver.LoadToRedis(rdb, "assets/bahrain_roads.json"); err != nil {
				log.Printf("failed to reload data: %v", err)
			} else {
				log.Println("Road data reloaded successfully.")
			}
		}
	}()

	router := gin.Default()
	routes.SetupRoutes(router, rdb)
	router.Run(":" + port)
}
