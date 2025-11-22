package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bsf-engine/receiver"

	"github.com/redis/go-redis/v9"
)

func GetSpeedLimitHandler(c *gin.Context, rdb *redis.Client) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing lat or lon query parameter"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
		return
	}

	speed, found := receiver.GetSpeedLimit(rdb, lat, lon)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no speed limit data found for location", "speed": 50})
		return
	}

	c.JSON(http.StatusOK, gin.H{"speed": speed})
}
