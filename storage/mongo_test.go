package storage

import (
	"context"
	"log"
	"testing"

	"github.com/adobromilskiy/pingatus/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	endpoint1 = &Endpoint{
		Name:   "test1",
		URL:    "http://test.com",
		Date:   1720688860,
		Status: true,
	}

	endpoint2 = &Endpoint{
		Name:   "test2",
		URL:    "http://test.com",
		Date:   1720688865,
		Status: true,
	}

	mongoURI = "mongodb://localhost:27117/testpingatus"
)

func TestGetMongoClient(t *testing.T) {
	cfg := &config.Config{
		MongoURI: mongoURI,
		Debug:    true,
	}

	store := GetMongoClient(cfg)
	if store == nil {
		t.Fatalf("Failed to get MongoDB client")
	}

	store.Close(context.Background())
}

func TestSaveEndpoint(t *testing.T) {
	ctx := context.Background()
	mongoopts := options.Client()

	mongoopts.ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, mongoopts)
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to MongoDB: %v", err)
	}

	store = &Store{
		Client: client,
		DBName: "testpingatus",
	}
	defer store.Close(ctx)

	err = store.SaveEndpoint(ctx, endpoint1)
	if err != nil {
		t.Fatalf("Failed to save endpoint 1: %v", err)
	}

	err = store.SaveEndpoint(ctx, endpoint1)
	if err != nil {
		t.Fatalf("Failed to save endpoint 2: %v", err)
	}

	err = store.SaveEndpoint(ctx, endpoint2)
	if err != nil {
		t.Fatalf("Failed to save endpoint 3: %v", err)
	}
}

func TestGetLastEndpoint(t *testing.T) {
	ctx := context.Background()
	mongoopts := options.Client()

	mongoopts.ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, mongoopts)
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to MongoDB: %v", err)
	}

	store = &Store{
		Client: client,
		DBName: "testpingatus",
	}
	defer store.Close(ctx)

	filter := bson.M{}
	endpoint, err := store.GetLastEndpoint(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to get last endpoint: %v", err)
	}

	if endpoint.Name != endpoint2.Name {
		t.Errorf("Last endpoint is unexpected: want %s, got %s", endpoint2.Name, endpoint.Name)
	}
}

func TestGetEndpoints(t *testing.T) {
	ctx := context.Background()
	mongoopts := options.Client()

	mongoopts.ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, mongoopts)
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to MongoDB: %v", err)
	}

	store = &Store{
		Client: client,
		DBName: "testpingatus",
	}
	defer store.Close(ctx)

	filter := bson.M{}
	endpoints, err := store.GetEndpoints(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to get endpoints: %v", err)
	}

	if len(endpoints) != 3 {
		t.Errorf("Endpoints count is unexpected: want 3, got %d", len(endpoints))
	}
}

func TestGetNames(t *testing.T) {
	ctx := context.Background()
	mongoopts := options.Client()

	mongoopts.ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, mongoopts)
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to MongoDB: %v", err)
	}

	store = &Store{
		Client: client,
		DBName: "testpingatus",
	}
	defer store.Close(ctx)

	names, err := store.GetNames(ctx)
	if err != nil {
		t.Fatalf("Failed to get endpoints: %v", err)
	}

	if len(names) != 2 {
		t.Errorf("Endpoints count is unexpected: want 2, got %d", len(names))
	}
}
