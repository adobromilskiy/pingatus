package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) getEndpoints(w http.ResponseWriter, r *http.Request) {
	data, err := s.db.GetEndpoints(r.Context())
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to get endpoints", "err", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)
	}

	err = WriteJSON(w, data)
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to write json", "err", err)
		http.Error(w, "failed to write json", http.StatusInternalServerError)
	}
}

func (s *Server) getEndpointStats(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)

		return
	}

	ago := time.Now().Unix() - 24*60*60

	endpoints, err := s.db.GetEndpointStats(r.Context(), name, ago)
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to get endpoint stats", "err", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)

		return
	}

	data := Stats{}
	data.Convert(endpoints)

	var total Duration
	for _, v := range endpoints {
		if v.Status {
			total++
		}
	}

	result := struct {
		Stats Stats  `json:"stats"`
		Total string `json:"total"`
	}{
		Stats: data,
		Total: total.String(),
	}

	err = WriteJSON(w, result)
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to write json", "err", err)
		http.Error(w, "failed to write json", http.StatusInternalServerError)
	}
}

func WriteJSON(w http.ResponseWriter, data any) error {
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(bts)))

	_, err = w.Write(bts)
	if err != nil {
		return err
	}

	return nil
}
