package entities

type Way struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	MaxSpeed string            `json:"maxspeed"`
	Tags     map[string]string `json:"tags"`
	Geometry [][2]float64      `json:"geometry"`
	MinLat   float64           `json:"-"`
	MinLon   float64           `json:"-"`
	MaxLat   float64           `json:"-"`
	MaxLon   float64           `json:"-"`
}

type Bounds struct {
	MinLat float64 `json:"minlat"`
	MinLon float64 `json:"minlon"`
	MaxLat float64 `json:"maxlat"`
	MaxLon float64 `json:"maxlon"`
}
