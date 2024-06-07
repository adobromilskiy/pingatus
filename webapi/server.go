package webapi

import (
	"context"
	"log"
	"net/http"

	"github.com/adobromilskiy/pingatus/config"
)

type Server struct {
	config.WEBAPIConfig
}

func NewServer(cfg config.WEBAPIConfig) *Server {
	return &Server{cfg}
}

func (s *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    s.ListenAddr,
		Handler: s.routes(),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] failed to shutdown webapi server: %v", err)
		} else {
			log.Printf("[INFO] webapi server stopped")
		}
	}()

	log.Printf("[INFO] webapi server listening on %s", s.ListenAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("[ERROR] failed to listen and serve: %v", err)
	}

	log.Printf("[INFO] webapi server stopped")
}

func (s *Server) routes() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	return r
}
