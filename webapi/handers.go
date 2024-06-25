package webapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) handlerGet24hrStats(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	ago := now - 24*60*60

	filter := bson.M{
		"name": name,
		"date": bson.M{
			"$gte": ago,
			"$lt":  now,
		},
	}

	endpoints, err := s.Store.GetEndpoints(r.Context(), filter)
	if err != nil {
		log.Println("[ERROR] failed to get endpoints:", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)
		return
	}

	data := Stats{}
	data.Convert(endpoints)

	var count Duration
	for _, v := range endpoints {
		if v.Status {
			count++
		}
	}

	result := struct {
		Stats  Stats  `json:"stats"`
		Uptime string `json:"uptime"`
	}{
		Stats:  data,
		Uptime: count.String(),
	}

	responce, err := json.Marshal(result)
	if err != nil {
		log.Println("[ERROR] failed to marshal response:", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responce)
}

func (s *Server) handlerGetCurrentStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	filter := bson.M{
		"name": name,
	}

	endpoint, err := s.Store.GetLastEndpoint(r.Context(), filter)
	if err != nil {
		log.Println("[ERROR] failed to get endpoints:", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)
		return
	}

	responce, err := json.Marshal(endpoint)
	if err != nil {
		log.Println("[ERROR] failed to marshal response:", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responce)
}
