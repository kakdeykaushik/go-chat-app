package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore[T any] struct {
	client *mongo.Client
	cfg    *Configuration
}

var singleMongoClient *mongo.Client

// singleton
func GetClient() (*mongo.Client, error) {
	if singleMongoClient != nil {
		return singleMongoClient, nil
	}

	uri := os.Getenv("mongo_uri")

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	singleMongoClient = c
	return singleMongoClient, err
}

func NewMongoStore[T any](db *mongo.Client, config *Configuration) *mongoStore[T] {
	return &mongoStore[T]{client: db, cfg: config}
}

func (ms *mongoStore[T]) Get(uid string) (*T, error) {
	var result T

	col := ms.client.Database(ms.cfg.DBName).Collection(ms.cfg.Collection)

	filter := bson.M{ms.cfg.Uid: uid}
	err := col.FindOne(context.Background(), filter).Decode(&result)

	return &result, err
}

func (ms *mongoStore[T]) List() ([]T, error) {
	panic("not implemented")
}

func (ms *mongoStore[T]) Save(data *T) error {
	col := ms.client.Database(ms.cfg.DBName).Collection(ms.cfg.Collection)

	_, err := col.InsertOne(context.Background(), &data)
	return err
}

func (ms *mongoStore[T]) Delete(id string) error {
	panic("not implemented")
}

func (ms *mongoStore[T]) Update(uid string, data *T) error {
	col := ms.client.Database(ms.cfg.DBName).Collection(ms.cfg.Collection)

	filter := bson.M{ms.cfg.Uid: uid}
	_, err := col.ReplaceOne(context.Background(), filter, &data)
	return err
}
