package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Endpoint struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	URL    string             `json:"url" bson:"url"`
	Status bool               `json:"status" bson:"status"`
	Date   int64              `json:"date" bson:"date"`
}
