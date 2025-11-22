package routes

import (
	"bsf-engine/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(router *gin.Engine, rdb *redis.Client) {
	router.GET("/get-speed-limit", func(c *gin.Context) {
		handlers.GetSpeedLimitHandler(c, rdb)
	})
}
