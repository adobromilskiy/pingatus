package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"gopkg.in/yaml.v3"
)

func TestGetMongoClient(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test get mongo client", func(mt *mtest.T) {
		tempFile, err := os.CreateTemp("", "config")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		os.Setenv("PINGATUS_CONFIG_PATH", tempFile.Name())

		expectedConfig := &config.Config{
			Debug:    true,
			MongoURI: "mongodb://localhost:27017/pingatus?timeoutMS=5000",
		}

		data, err := yaml.Marshal(expectedConfig)
		if err != nil {
			t.Fatalf("Failed to marshal expected config: %v", err)
		}

		if _, err := tempFile.Write(data); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		store := GetMongoClient()

		assert.NotNil(t, store)
		assert.NotNil(t, store.Client)
		assert.Equal(t, "pingatus", store.DBName)

		store.Close()
	})
}

func TestSaveEndpoint(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test save endpoint", func(mt *mtest.T) {
		store := &Store{
			Client: mt.Client,
			DBName: mt.DB.Name(),
		}

		endpoint := &Endpoint{
			Name:   "test",
			URL:    "http://test.com",
			Date:   time.Now().Unix(),
			Status: true,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := store.SaveEndpoint(context.TODO(), endpoint)

		assert.Nil(t, err)
	})

}

func TestGetLastEndpoint(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find one success", func(mt *mtest.T) {
		store := &Store{
			Client: mt.Client,
			DBName: mt.DB.Name(),
		}

		filter := bson.M{}
		expectedEndpoint := &Endpoint{
			ID:   primitive.NewObjectID(),
			Name: "Last Endpoint",
			URL:  "http://last-endpoint.com",
		}

		bsonBytes, err := bson.Marshal(expectedEndpoint)
		if err != nil {
			t.Fatalf("Failed to marshal endpoint: %v", err)
		}

		var doc primitive.D
		err = bson.Unmarshal(bsonBytes, &doc)
		if err != nil {
			t.Fatalf("Failed to unmarshal BSON into primitive.D: %v", err)
		}

		firstBatch := mtest.CreateCursorResponse(1, "test.endpoints", mtest.FirstBatch, doc)
		mt.AddMockResponses(firstBatch)

		endpoint, err := store.GetLastEndpoint(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedEndpoint.ID, endpoint.ID)
		assert.Equal(t, expectedEndpoint.Name, endpoint.Name)
		assert.Equal(t, expectedEndpoint.URL, endpoint.URL)

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "mock error",
		}))

		endpoint, err = store.GetLastEndpoint(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, endpoint)
	})
}

func TestGetLastEndpoint_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find error", func(mt *mtest.T) {
		store := &Store{
			Client: mt.Client,
			DBName: mt.DB.Name(),
		}

		filter := bson.M{}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "mock error",
		}))

		endpoint, err := store.GetLastEndpoint(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, endpoint)
	})
}

func TestGetEndpoints(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find all success", func(mt *mtest.T) {
		store := &Store{
			Client: mt.Client,
			DBName: mt.DB.Name(),
		}

		filter := bson.M{}
		expectedEndpoints := []*Endpoint{
			{
				ID:     primitive.NewObjectID(),
				Name:   "test",
				URL:    "http://example1.com",
				Status: true,
				Date:   time.Now().Unix(),
			},
			{
				ID:     primitive.NewObjectID(),
				Name:   "test",
				URL:    "http://example2.com",
				Status: true,
				Date:   time.Now().Unix(),
			},
		}

		bsonBytes, err := bson.Marshal(expectedEndpoints[0])
		if err != nil {
			t.Fatalf("Failed to marshal endpoint: %v", err)
		}

		var doc primitive.D
		err = bson.Unmarshal(bsonBytes, &doc)
		if err != nil {
			t.Fatalf("Failed to unmarshal BSON into primitive.D: %v", err)
		}

		bsonBytes, err = bson.Marshal(expectedEndpoints[0])
		if err != nil {
			t.Fatalf("Failed to marshal endpoint: %v", err)
		}

		var doc2 primitive.D
		err = bson.Unmarshal(bsonBytes, &doc)
		if err != nil {
			t.Fatalf("Failed to unmarshal BSON into primitive.D: %v", err)
		}

		first := mtest.CreateCursorResponse(1, "test.endpoints", mtest.FirstBatch, doc)
		second := mtest.CreateCursorResponse(1, "test.endpoints", mtest.NextBatch, doc2)
		killCursors := mtest.CreateCursorResponse(0, "test.endpoints", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		endpoints, err := store.GetEndpoints(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedEndpoints[0].Date, endpoints[0].Date)
	})
}

func TestGetEndpoints_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find all err", func(mt *mtest.T) {
		store := &Store{
			Client: mt.Client,
			DBName: mt.DB.Name(),
		}

		filter := bson.M{}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "mock error",
		}))

		endpoints, err := store.GetEndpoints(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, endpoints)
	})
}
