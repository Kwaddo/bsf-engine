package receiver

import (
	"encoding/json"
	"os"

	"bsf-engine/entities"

	"github.com/redis/go-redis/v9"
)

const pipelineBatch = 1000

func LoadToRedis(rdb *redis.Client, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pipe := rdb.Pipeline()
	count := 0

	var ways []entities.Way
	if err := json.NewDecoder(file).Decode(&ways); err != nil {
		return err
	}

	for i := range ways {
		minLat, minLon := 999.0, 999.0
		maxLat, maxLon := -999.0, -999.0
		for _, pt := range ways[i].Geometry {
			lat, lon := pt[0], pt[1]
			if lat < minLat {
				minLat = lat
			}
			if lat > maxLat {
				maxLat = lat
			}
			if lon < minLon {
				minLon = lon
			}
			if lon > maxLon {
				maxLon = lon
			}
		}
		ways[i].MinLat = minLat
		ways[i].MinLon = minLon
		ways[i].MaxLat = maxLat
		ways[i].MaxLon = maxLon

		if err := EnqueueWay(pipe, &ways[i]); err != nil {
			return err
		}
		count++
		if count%pipelineBatch == 0 {
			if _, err := pipe.Exec(entities.CTX); err != nil {
				return err
			}
			pipe = rdb.Pipeline()
		}

		// Populate in-memory cache
		speed, _ := entities.ParseSpeed(ways[i].MaxSpeed)
		waysCache.Lock()
		waysCache.byID[ways[i].ID] = CachedWay{
			Geometry: ways[i].Geometry,
			Speed:    speed,
		}
		waysCache.Unlock()
	}
	if _, err := pipe.Exec(entities.CTX); err != nil {
		return err
	}
	return nil
}
