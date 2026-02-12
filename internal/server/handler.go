package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var (
	errStatsRangeMissing = errors.New("from and to are required")
	errStatsFromInvalid  = errors.New("from is invalid")
	errStatsToInvalid    = errors.New("to is invalid")
	errStatsRangeInvalid = errors.New("from must be <= to")
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

	from, to, err := parseStatsRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	endpoints, err := s.db.GetEndpointStats(r.Context(), name, from, to)
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

func (s *Server) getEndpointByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "name is required", http.StatusBadRequest)

		return
	}

	from, to, err := parseRequiredStatsRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	endpoints, err := s.db.GetEndpointStats(r.Context(), name, from, to)
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to get endpoint stats", "err", err)
		http.Error(w, "failed to get endpoints", http.StatusInternalServerError)

		return
	}

	err = WriteJSON(w, endpoints)
	if err != nil {
		s.lg.Log(r.Context(), slog.LevelError, "failed to write json", "err", err)
		http.Error(w, "failed to write json", http.StatusInternalServerError)
	}
}

func parseStatsRange(r *http.Request) (int64, int64, error) {
	fromRaw := r.URL.Query().Get("from")
	toRaw := r.URL.Query().Get("to")

	if fromRaw == "" && toRaw == "" {
		const daySeconds = 24 * 60 * 60

		to := time.Now().Unix()
		from := to - daySeconds

		return from, to, nil
	}

	if fromRaw == "" || toRaw == "" {
		return 0, 0, errStatsRangeMissing
	}

	from, err := parseTimeParam(fromRaw)
	if err != nil {
		return 0, 0, errStatsFromInvalid
	}

	to, err := parseTimeParam(toRaw)
	if err != nil {
		return 0, 0, errStatsToInvalid
	}

	if from > to {
		return 0, 0, errStatsRangeInvalid
	}

	return from, to, nil
}

func parseRequiredStatsRange(r *http.Request) (int64, int64, error) {
	fromRaw := r.URL.Query().Get("from")
	toRaw := r.URL.Query().Get("to")

	if fromRaw == "" || toRaw == "" {
		return 0, 0, errStatsRangeMissing
	}

	from, err := parseTimeParam(fromRaw)
	if err != nil {
		return 0, 0, errStatsFromInvalid
	}

	to, err := parseTimeParam(toRaw)
	if err != nil {
		return 0, 0, errStatsToInvalid
	}

	if from > to {
		return 0, 0, errStatsRangeInvalid
	}

	return from, to, nil
}

func parseTimeParam(raw string) (int64, error) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return 0, err
	}

	return parsed.UTC().Unix(), nil
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
