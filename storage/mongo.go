package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var (
	mongoClient *mongo.Client
	mongoError  error
	mongoOnce   sync.Once
)

func GetMongoClient() (*mongo.Client, error) {
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

		fmt.Println("Connected to MongoDB!")

		client.Database(cs.Database)

		mongoClient = client
	})

	return mongoClient, mongoError
}
