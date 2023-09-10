package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"chat-app/pkg/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type roomSvc struct {
	// DB domain.Storage
}

func NewRoomSvc() roomSvc {
	// return roomSvc{DB: store}
	return roomSvc{}
}

func getCollection(dbName, collectioName string) *mongo.Collection {
	client := db.GetClient()
	col := client.Database(dbName).Collection(collectioName)
	return col
}

func (rs roomSvc) GetRoom(roomId string) (*model.Room, error) {

	col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_ROOM)
	filter := bson.M{"roomId": roomId}

	var result *entity.Room
	err := col.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Error while getting room", err)
		return nil, err
	}

	room := utils.EntityToModelRoom(result)

	return room, nil
}

func (rs roomSvc) CreateRoom() (*model.Room, error) {
	// roomId := uuid.New().String()[:5]
	roomId := "abcde"
	room := &model.Room{RoomId: roomId, Members: []*model.Member{}}

	col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_ROOM)

	roomEntity := utils.ModelToEntityRoom(room)
	_, err := col.InsertOne(context.TODO(), roomEntity)
	if err != nil {
		fmt.Println("Error while creating room", err)
		return nil, err
	}
	return room, err
}

func (rs roomSvc) AddMember(room *model.Room, member *model.Member) {
	room.Members = append(room.Members, member)

	col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_ROOM)
	filter := bson.M{"roomId": room.RoomId}

	roomEntity := utils.ModelToEntityRoom(room)
	_, err := col.ReplaceOne(context.TODO(), filter, &roomEntity)
	if err != nil {
		fmt.Println("Error while adding member to the room", err)
	}
}

func (rs roomSvc) RemoveMember(room *model.Room, username string) error {
	// room.Lock()
	// defer room.Unlock()

	for i, member := range room.Members {
		if member.Username == username {
			room.Members = utils.RemoveIndex(room.Members, i)
			// update db
			col := getCollection(utils.DB_CHATROOM, utils.COLLECTION_ROOM)
			filter := bson.M{"roomId": room.RoomId}

			roomEntity := utils.ModelToEntityRoom(room)

			_, err := col.ReplaceOne(context.TODO(), filter, &roomEntity)
			if err != nil {
				fmt.Println("Error while removing member", err)
				return err
			}

			break
		}
	}

	return nil
}

func (rs roomSvc) IsNewMember(room *model.Room, memberUsername string) bool {
	for _, member := range room.Members {
		if member.Username == memberUsername {
			return false
		}
	}
	return true
}
