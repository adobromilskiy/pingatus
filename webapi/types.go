package webapi

import (
	"time"

	"github.com/adobromilskiy/pingatus/storage"
)

type (
	Stats struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		// Hours  []int  `json:"points"`
		// Points []int  `json:"points"`
		Points map[int]int `json:"points"`
	}
)

func (s *Stats) Convert(endpoints []*storage.Endpoint) {
	if len(endpoints) == 0 {
		return
	}

	stats := make(map[int]int)
	points := 0
	counts := 0
	checkhour := -1

	for _, endpoint := range endpoints {
		hour := time.Unix(endpoint.Date, 0).Hour()
		if checkhour != hour {
			if checkhour != -1 {
				stats[checkhour] = points * 100 / counts
			}
			checkhour = hour
			points = 0
			counts = 0
		}
		counts++
		if endpoint.Status {
			points++
		}

		if endpoint == endpoints[len(endpoints)-1] {
			stats[checkhour] = points * 100 / counts
		}
	}

	s.Name = endpoints[0].Name
	s.URL = endpoints[0].URL
	s.Points = stats

}
