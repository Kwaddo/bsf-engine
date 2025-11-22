package receiver

import (
	"encoding/json"
	"fmt"

	"bsf-engine/entities"

	"github.com/redis/go-redis/v9"
)

func EnqueueWay(pipe redis.Pipeliner, w *entities.Way) error {
	speed, _ := entities.ParseSpeed(w.MaxSpeed)
	if speed == 0 {
		speed = 50
	}

	geometryJSON, err := json.Marshal(w.Geometry)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("way:%d", w.ID)
	pipe.HSet(entities.CTX, key, map[string]interface{}{
		"geometry": string(geometryJSON),
		"maxspeed": speed,
		"name":     w.Name,
	})

	if len(w.Geometry) == 0 {
		return nil
	}

	const chunk = 200
	locations := make([]*redis.GeoLocation, 0, chunk)
	flushLocations := func() {
		if len(locations) == 0 {
			return
		}
		pipe.GeoAdd(entities.CTX, "ways:index", locations...)
		locations = locations[:0]
	}

	for idx, pt := range w.Geometry {
		locations = append(locations, &redis.GeoLocation{
			Name:      fmt.Sprintf("%d:%d", w.ID, idx),
			Longitude: pt[1],
			Latitude:  pt[0],
		})
		if len(locations) == chunk {
			flushLocations()
		}
	}
	flushLocations()
	return nil
}
