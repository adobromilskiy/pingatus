package webapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/storage"
	"go.mongodb.org/mongo-driver/bson"
)

func HandlerGet24hrStats(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	store := storage.GetMongoClient()

	now := time.Now().Unix()
	ago := now - 24*60*60

	filter := bson.M{
		"name": name,
		"date": bson.M{
			"$gte": ago,
			"$lt":  now,
		},
	}

	endpoints, err := store.GetEndpoints(r.Context(), filter)
	if err != nil {
		log.Println("[ERROR] failed to get endpoints:", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)
		return
	}

	data := Stats{}
	data.Convert(endpoints)

	responce, err := json.Marshal(data)
	if err != nil {
		log.Println("[ERROR] failed to marshal response:", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responce)
}

func HandlerGetCurrentStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	store := storage.GetMongoClient()

	filter := bson.M{
		"name": name,
	}

	endpoint, err := store.GetLastEndpoint(r.Context(), filter)
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
