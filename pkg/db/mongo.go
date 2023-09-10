package db

import (
	"chat-app/pkg/types"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	DBName     string
	Collection string
}

type mongoStore struct {
	Client *mongo.Client
	Config *Configuration
}

var client *mongo.Client

func GetClient() *mongo.Client {
	if client != nil {
		return client
	}

	uri := os.Getenv("mongo_uri")
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil
	}

	client = c
	return client
}

func newMongoStore(config *Configuration) types.Storage {
	client := GetClient()
	return &mongoStore{Client: client, Config: config}
}

func (ms *mongoStore) Get(id string) (any, error) {
	panic("not implemented")
}

func (ms *mongoStore) List() ([]any, error) {
	panic("not implemented")
}

func (ms *mongoStore) Save(K, data any) error {
	panic("not implemented")
}

func (ms *mongoStore) Delete(id string) error {
	panic("not implemented")
}
