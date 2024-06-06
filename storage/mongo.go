package storage

import (
	"context"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store struct {
	Client *mongo.Client
}

var (
	store      *Store
	mongoError error
	mongoOnce  sync.Once
)

func GetMongoClient() (*Store, error) {
	mongoOnce.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			mongoError = err
			return
		}
		cs, err := connstring.ParseAndValidate(cfg.MongoURI)
		if err != nil {
			mongoError = err
			return
		}

		mongoopts := options.Client()
		mongoopts.ApplyURI(cfg.MongoURI)

		if cfg.Debug {
			monitor := &event.CommandMonitor{
				Started: func(_ context.Context, event *event.CommandStartedEvent) {
					log.Println("[DEBUG] Started:", event.CommandName, event.RequestID, event.ConnectionID, event.Command)
				},
				Succeeded: func(_ context.Context, event *event.CommandSucceededEvent) {
					log.Println("[DEBUG] Succeeded:", event.CommandName, event.RequestID, event.ConnectionID)
				},
				Failed: func(_ context.Context, event *event.CommandFailedEvent) {
					log.Println("[DEBUG] Failed:", event.CommandName, event.RequestID, event.ConnectionID)
				},
			}
			mongoopts.SetMonitor(monitor)
		}

		ctx := context.Background()

		client, err := mongo.Connect(ctx, mongoopts)
		if err != nil {
			mongoError = err
			return
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			mongoError = err
			return
		}

		client.Database(cs.Database)

		store = new(Store)
		store.Client = client

		log.Println("[INFO] connected to MongoDB!")
	})

	return store, mongoError
}

func (s *Store) Close() {
	if s.Client != nil {
		if err := s.Client.Disconnect(context.Background()); err != nil {
			log.Println("[ERROR] disconnecting from MongoDB:", err)
		}
		log.Println("[INFO] disconnected from MongoDB!")
	}
}
