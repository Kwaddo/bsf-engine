package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OverpassResponse struct {
	Elements []struct {
		ID    int64             `json:"id"`
		Type  string            `json:"type"`
		Tags  map[string]string `json:"tags"`
		Nodes []int64           `json:"nodes"`
		Lat   float64           `json:"lat,omitempty"`
		Lon   float64           `json:"lon,omitempty"`
	} `json:"elements"`
}

type RoadInfo struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	MaxSpeed string            `json:"maxspeed"`
	Tags     map[string]string `json:"tags"`
	Geometry [][2]float64      `json:"geometry"`
}

func main() {
	url := "https://overpass-api.de/api/interpreter"
	query := `
[out:json][timeout:180];
(
  way["highway"](25.796,50.269,26.366,50.720);
  node(w);
);
out body;
`

	resp, err := http.Post(url, "application/x-www-form-urlencoded",
		io.NopCloser(strings.NewReader("data="+query)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Failed to decode JSON. Response was:")
		fmt.Println(string(body))
		panic(err)
	}

	nodeCoords := make(map[int64][2]float64)
	for _, el := range data.Elements {
		if el.Type == "node" {
			nodeCoords[el.ID] = [2]float64{el.Lat, el.Lon}
		}
	}

	var roads []RoadInfo
	for _, el := range data.Elements {
		if el.Type != "way" {
			continue
		}
		var geometry [][2]float64
		for _, nid := range el.Nodes {
			if coord, ok := nodeCoords[nid]; ok {
				geometry = append(geometry, coord)
			}
		}
		name := el.Tags["name"]
		if name == "" {
			name = fmt.Sprintf("unknown_%d", el.ID)
		}
		maxspeed := el.Tags["maxspeed"]
		if maxspeed == "" {
			maxspeed = "50"
		}
		roads = append(roads, RoadInfo{
			ID:       el.ID,
			Name:     name,
			MaxSpeed: maxspeed,
			Tags:     el.Tags,
			Geometry: geometry,
		})
	}

	f, err := os.Create("../assets/bahrain_roads.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(roads)

	fmt.Printf("Total roads: %d\n", len(roads))
}
