package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

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
