package storage

import (
	"context"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store struct {
	Client *mongo.Client
	DBName string
}

var (
	store     *Store
	mongoOnce sync.Once
)

func GetMongoClient(cfg *config.Config) Storage {
	mongoOnce.Do(func() {
		cs, err := connstring.ParseAndValidate(cfg.MongoURI)
		if err != nil {
			log.Fatalf("[ERROR] failed to parse MongoDB URI: %v", err)
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
			log.Fatalf("[ERROR] failed to connect to MongoDB: %v", err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("[ERROR] failed to ping MongoDB: %v", err)
		}

		store = &Store{
			Client: client,
			DBName: cs.Database,
		}

		log.Println("[INFO] connected to MongoDB!")
	})

	return store
}

func (s *Store) Close() {
	if s.Client != nil {
		if err := s.Client.Disconnect(context.Background()); err != nil {
			log.Println("[ERROR] disconnecting from MongoDB:", err)
		}
		log.Println("[INFO] disconnected from MongoDB!")
	}
}

func (s *Store) SaveEndpoint(ctx context.Context, endpoint *Endpoint) error {
	collection := s.Client.Database(s.DBName).Collection("endpoints")
	_, err := collection.InsertOne(ctx, endpoint)
	return err
}

func (s *Store) GetEndpoints(ctx context.Context, filter primitive.M) ([]*Endpoint, error) {
	collection := s.Client.Database(s.DBName).Collection("endpoints")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var endpoints []*Endpoint
	for cursor.Next(ctx) {
		var endpoint Endpoint
		if err := cursor.Decode(&endpoint); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, nil
}

func (s *Store) GetLastEndpoint(ctx context.Context, filter primitive.M) (*Endpoint, error) {
	collection := s.Client.Database(s.DBName).Collection("endpoints")
	opts := options.FindOne().SetSort(bson.M{"date": -1})
	var endpoint Endpoint
	err := collection.FindOne(ctx, filter, opts).Decode(&endpoint)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}
