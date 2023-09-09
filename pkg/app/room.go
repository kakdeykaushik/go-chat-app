package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/domain"
	"chat-app/shared"
	"context"

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

func (rs roomSvc) GetRoom(roomId string) (*domain.Room, error) {

	col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_ROOM)
	filter := bson.M{"roomId": roomId}

	var result *domain.Room
	err := col.FindOne(context.TODO(), filter).Decode(&result)
	shared.HandleError(err, "Error while getting room")

	return result, nil
}

func (rs roomSvc) CreateRoom() (*domain.Room, error) {
	// roomId := uuid.New().String()[:5]
	roomId := "abcde"
	room := &domain.Room{RoomId: roomId, Members: nil}

	col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_ROOM)

	_, err := col.InsertOne(context.TODO(), room)
	shared.HandleError(err, "Error while creating room")
	return room, err
}

func (rs roomSvc) AddMember(room *domain.Room, member *domain.Member) {
	room.Members = append(room.Members, member)

	col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_ROOM)
	filter := bson.M{"roomId": room.RoomId}

	_, err := col.ReplaceOne(context.TODO(), filter, &room)
	shared.HandleError(err, "Error while adding member to the room")
}

func (rs roomSvc) RemoveMember(room *domain.Room, username string) {
	// room.Lock()
	// defer room.Unlock()

	for i, member := range room.Members {
		if member.Username == username {
			room.Members = shared.RemoveIndex(room.Members, i)
			// update db
			col := getCollection(shared.DB_CHATROOM, shared.COLLECTION_ROOM)
			filter := bson.M{"roomId": room.RoomId}

			_, err := col.ReplaceOne(context.TODO(), filter, &room)
			shared.HandleError(err, "Error while removing member")

			break
		}
	}

}

func (rs roomSvc) IsNewMember(room *domain.Room, memberUsername string) bool {
	for _, member := range room.Members {
		if member.Username == memberUsername {
			return false
		}
	}
	return true
}
