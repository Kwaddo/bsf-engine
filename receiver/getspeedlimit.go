package receiver

import (
	"math"
	"strconv"
	"strings"

	"bsf-engine/entities"

	"github.com/redis/go-redis/v9"
)

func GetSpeedLimit(rdb *redis.Client, lat, lon float64) (int, bool) {
	radii := []float64{400, 800, 1500}
	for _, radius := range radii {
		res, err := rdb.GeoRadius(entities.CTX, "ways:index", lon, lat, &redis.GeoRadiusQuery{
			Radius:    radius,
			Unit:      "m",
			WithCoord: false,
			WithDist:  true,
			Sort:      "ASC",
			Count:     80,
		}).Result()
		if err != nil {
			return 50, false
		}
		_, bestSpeed, found := snapToNearest(res, lat, lon)
		if found {
			return bestSpeed, true
		}
	}

	return 50, false
}

func snapToNearest(candidates []redis.GeoLocation, lat, lon float64) (float64, int, bool) {
	const (
		primarySnap  = 40.0
		fallbackSnap = 90.0
	)
	bestDist := math.MaxFloat64
	bestSpeed := 50

	for _, candidate := range candidates {
		parts := strings.SplitN(candidate.Name, ":", 2)
		if len(parts) == 0 {
			continue
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		waysCache.RLock()
		cw, ok := waysCache.byID[id]
		waysCache.RUnlock()
		if !ok || len(cw.Geometry) == 0 {
			continue
		}
		dist := entities.PointToLineStringDistance(lat, lon, cw.Geometry)
		if dist < bestDist {
			bestDist = dist
			bestSpeed = cw.Speed
			if dist <= primarySnap {
				break
			}
		}
	}

	if bestDist <= fallbackSnap {
		return bestDist, bestSpeed, true
	}
	return bestDist, bestSpeed, false
}
