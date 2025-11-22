package entities

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

func (b Bounds) Contains(lat, lon float64) bool {
	return lat >= b.MinLat && lat <= b.MaxLat &&
		lon >= b.MinLon && lon <= b.MaxLon
}

func PointToLineStringDistance(lat, lon float64, lineString [][2]float64) float64 {
	if len(lineString) == 0 {
		return math.MaxFloat64
	}

	minDist := math.MaxFloat64
	for i := 0; i < len(lineString)-1; i++ {
		dist := PointToSegmentDistance(lat, lon, lineString[i], lineString[i+1])
		if dist < minDist {
			minDist = dist
		}
	}
	return minDist
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	dφ := (lat2 - lat1) * math.Pi / 180
	dλ := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dφ/2)*math.Sin(dφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(dλ/2)*math.Sin(dλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func PointToSegmentDistance(lat, lon float64, p1, p2 [2]float64) float64 {
	lat1, lon1 := p1[0], p1[1]
	lat2, lon2 := p2[0], p2[1]

	x0, y0 := lon, lat
	x1, y1 := lon1, lat1
	x2, y2 := lon2, lat2

	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 && dy == 0 {
		return Haversine(lat, lon, lat1, lon1)
	}

	t := ((x0-x1)*dx + (y0-y1)*dy) / (dx*dx + dy*dy)
	if t < 0 {
		return Haversine(lat, lon, lat1, lon1)
	} else if t > 1 {
		return Haversine(lat, lon, lat2, lon2)
	}
	projLon := x1 + t*dx
	projLat := y1 + t*dy
	return Haversine(lat, lon, projLat, projLon)
}

func ParseSpeed(raw interface{}) (int, bool) {
	var s string
	switch v := raw.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return 0, false
	}
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" || s == "signals" {
		return 0, false
	}

	multiplier := 1.0
	switch {
	case strings.HasSuffix(s, "mph"):
		multiplier = 1.60934
		s = strings.TrimSpace(strings.TrimSuffix(s, "mph"))
	case strings.HasSuffix(s, "km/h"):
		s = strings.TrimSpace(strings.TrimSuffix(s, "km/h"))
	}

	digits := regexp.MustCompile(`\d+`).FindString(s)
	if digits == "" {
		return 0, false
	}
	val, err := strconv.Atoi(digits)
	if err != nil {
		return 0, false
	}
	return int(math.Round(float64(val) * multiplier)), true
}
