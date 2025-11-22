package receiver

import (
	"sync"
)

type CachedWay struct {
	Geometry [][2]float64
	Speed    int
}

var waysCache = struct {
	sync.RWMutex
	byID map[int]CachedWay
}{byID: make(map[int]CachedWay)}
