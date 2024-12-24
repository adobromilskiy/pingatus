package server

import (
	"fmt"
	"time"

	"github.com/adobromilskiy/pingatus/core"
)

type (
	Stats struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Hours  []int  `json:"hours"`
		Points []int  `json:"points"`
	}

	Duration int
)

func (s *Stats) Convert(endpoints []core.Endpoint) {
	if len(endpoints) == 0 {
		return
	}

	points := 0
	counts := 0
	checkhour := -1

	for _, endpoint := range endpoints {
		hour := time.Unix(endpoint.Date, 0).Hour()
		if checkhour != hour {
			if checkhour != -1 {
				s.Hours = append(s.Hours, checkhour)
				s.Points = append(s.Points, points*100/counts)
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
			s.Hours = append(s.Hours, checkhour)
			s.Points = append(s.Points, points*100/counts)
		}
	}

	s.Name = endpoints[0].Name
	s.URL = endpoints[0].Address
}

func (d Duration) String() string {
	hour := d / 60
	minute := d % 60

	var result string
	result = fmt.Sprintf("%d", hour)

	if hour < 10 {
		result = fmt.Sprintf("0%d", hour)
	}

	if minute < 10 {
		result = fmt.Sprintf("%s:0%d", result, minute)
	} else {
		result = fmt.Sprintf("%s:%d", result, minute)
	}

	return result
}
