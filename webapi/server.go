package webapi

import (
	"context"
	"log"
	"net/http"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/adobromilskiy/pingatus/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	config.WEBAPIConfig
	Store Storage
}

type Storage interface {
	GetLastEndpoint(ctx context.Context, filter primitive.M) (*storage.Endpoint, error)
	GetEndpoints(ctx context.Context, filter primitive.M) ([]*storage.Endpoint, error)
	SaveEndpoint(ctx context.Context, endpoint *storage.Endpoint) error
	Close()
}

func NewServer(cfg config.WEBAPIConfig, store Storage) *Server {
	return &Server{cfg, store}
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
	r.Handle("/", http.FileServer(http.Dir(s.AssetsDir)))
	r.HandleFunc("/api/24hrstats", s.handlerGet24hrStats)
	r.HandleFunc("/api/currentstatus", s.handlerGetCurrentStatus)
	return r
}
