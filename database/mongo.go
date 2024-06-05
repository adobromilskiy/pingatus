package database

import (
	"context"
	"fmt"
	"log"
	"sync"

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

func GetMongoClient(ctx context.Context, uri string, debug bool) (*mongo.Client, error) {
	mongoOnce.Do(func() {
		cs, err := connstring.ParseAndValidate(uri)
		if err != nil {
			mongoError = err
			return
		}

		mongoopts := options.Client()
		mongoopts.ApplyURI(uri)

		if debug {
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
