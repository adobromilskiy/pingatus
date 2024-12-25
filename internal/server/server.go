package server

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/adobromilskiy/pingatus/core"
)

//go:embed *.html
var content embed.FS

type Server struct {
	lg   *slog.Logger
	db   core.Setter
	addr string
}

func New(lg *slog.Logger, db core.Setter, addr string) *Server {
	return &Server{lg.With("pkg", "server"), db, addr}
}

func (s *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:              s.addr,
		Handler:           s.routes(),
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(ctx); err != nil {
			s.lg.Log(ctx, slog.LevelError, "failed to shutdown server", "err", err)
		}

		s.lg.Log(ctx, slog.LevelInfo, "server stopped")
	}()

	s.lg.Log(ctx, slog.LevelInfo, "server started", "addr", s.addr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.lg.Log(ctx, slog.LevelError, "failed to listen and serve", "err", err)
	}

	s.lg.Log(ctx, slog.LevelInfo, "server stopped")
}

func (s *Server) routes() http.Handler {
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.FS(content)))
	r.HandleFunc("/stats", s.getEndpointStats)
	r.HandleFunc("/endpoints", s.getEndpoints)

	return r
}
