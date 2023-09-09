package db

import (
	"chat-app/pkg/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	// xyzRepository repository.XyzRepository
	// this can be used to plug different db altogther, exactly what I... wanted ?
}

var client *mongo.Client

func newMongoStore(config *Configuration) domain.Storage {
	// uri := os.Getenv("mongo_uri")
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return nil
	}
	return mongoStore{Client: client, Config: config}
	// return mongoStore[domain.Room]{Client: client, Config: config, result: domain.Room{}}
}

func GetClient() *mongo.Client {

	if client != nil {
		return client
	}

	uri := "mongodb://localhost:27017"
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil
	}

	client = c

	return client

}

func (ms mongoStore) Get(id string) (any, error) {

	return nil, nil
	// col := ms.Client.Database(ms.Config.DBName).Collection(ms.Config.Collection)
	// // retu
	// filter := bson.M{"roomId": id}

	// var result *D

	// switch result.(type) {
	// case domain.Room:
	// 	var result *domain.Room
	// 	_ = col.FindOne(context.TODO(), filter).Decode(result)
	// 	return result, nil
	// case domain.Member:
	// 	var result *domain.Member
	// 	_ = col.FindOne(context.TODO(), filter).Decode(result)
	// 	return result, nil
	// }

	// fmt.Printf("mango1 - %+v %T\n", result, result)
	// err := col.FindOne(context.TODO(), filter).Decode(&result)
	// fmt.Printf("mango2 - %+v %T\n", result, result)

	// if err != nil {
	// 	return nil, err
	// }
	// return result, nil
}

func (ms mongoStore) List() ([]any, error) {
	col := ms.Client.Database(ms.Config.DBName).Collection(ms.Config.Collection)
	filter := bson.M{}
	result := []any{}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ms mongoStore) Save(K, data any) error {
	col := ms.Client.Database(ms.Config.DBName).Collection(ms.Config.Collection)
	_, err := col.InsertOne(context.TODO(), data)
	return err
}

func (ms mongoStore) Delete(id string) error {
	col := ms.Client.Database(ms.Config.DBName).Collection(ms.Config.Collection)
	filter := bson.M{"id": id}
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}
