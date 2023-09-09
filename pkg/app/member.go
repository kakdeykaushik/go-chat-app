package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/domain"
	"chat-app/shared"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type memberSvc struct {
	DB domain.Storage
}

func NewMemberSvc(store domain.Storage) memberSvc {
	return memberSvc{DB: store}
}

func (ms memberSvc) CreateMember(username string, socket *websocket.Conn) *domain.Member {
	conn := domain.Connection{Socket: socket, Mutex: sync.Mutex{}}
	member := &domain.Member{Username: username, Conn: &conn}

	client := db.GetClient()
	col := client.Database("chatroom").Collection("member")

	_, err := col.InsertOne(context.TODO(), member)

	shared.HandleError(err, "Error while creating member")

	return member
}

func (ms memberSvc) GetMember(username string) (*domain.Member, error) {

	client := db.GetClient()
	col := client.Database("chatroom").Collection("member")

	filter := bson.M{"username": username}

	var member *domain.Member
	err := col.FindOne(context.Background(), filter).Decode(&member)

	return member, err
}
