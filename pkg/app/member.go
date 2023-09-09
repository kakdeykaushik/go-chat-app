package app

import (
	"chat-app/pkg/domain"
	"chat-app/shared"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberSvc struct {
	// DB domain.Storage
}

func NewMemberSvc() memberSvc {
	// return memberSvc{DB: store}
	return memberSvc{}
}

func (ms memberSvc) CreateMember(username string, socket *websocket.Conn) *domain.Member {

	m, err := ms.GetMember(username)
	if err == mongo.ErrNoDocuments {
		conn := domain.Connection{Socket: socket, Mutex: sync.Mutex{}}
		member := &domain.Member{Username: username, Conn: &conn}

		col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_MEMBER)
		_, err := col.InsertOne(context.TODO(), member)

		shared.HandleError(err, "Error while creating member")
		return member
	}

	shared.HandleError(err, "Error creating member")
	return m
}

func (ms memberSvc) GetMember(username string) (*domain.Member, error) {
	col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_MEMBER)

	var member *domain.Member

	filter := bson.M{"username": username}
	err := col.FindOne(context.Background(), filter).Decode(&member)

	return member, err
}
