package app

import (
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"chat-app/pkg/utils"
	"context"
	"fmt"
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

func (ms memberSvc) CreateMember(username string, socket *websocket.Conn) *model.Member {

	m, err := ms.GetMember(username)
	if err == mongo.ErrNoDocuments {

		conn := model.Connection{Socket: socket, Mutex: &sync.Mutex{}}
		member := &model.Member{Username: username, Conn: &conn}

		col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_MEMBER)
		memberEntity := utils.ModelToEntityMember(member)
		_, err := col.InsertOne(context.TODO(), memberEntity)

		if err != nil {
			fmt.Println("Error while creating member", err)
			return nil
		}
		return member
	}

	if err != nil {
		fmt.Println("Error creating member", err)
		return nil
	}

	return m
}

func (ms memberSvc) GetMember(username string) (*model.Member, error) {
	col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_MEMBER)

	var memberEntity *entity.Member

	filter := bson.M{"username": username}
	err := col.FindOne(context.Background(), filter).Decode(&memberEntity)

	if err == mongo.ErrNoDocuments {
		return nil, err
	}

	member := utils.EntityToModelMember(memberEntity)
	return member, err
}
