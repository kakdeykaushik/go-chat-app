package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/domain"
	"chat-app/shared"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type roomSvc struct {
	DB domain.Storage
}

func NewRoomSvc(store domain.Storage) roomSvc {
	return roomSvc{DB: store}
}

func (rs roomSvc) GetRoom(roomId string) (*domain.Room, error) {

	client := db.GetClient()

	col := client.Database("chatroom").Collection("room")
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

	client := db.GetClient()

	col := client.Database("chatroom").Collection("room")
	_, err := col.InsertOne(context.TODO(), room)
	shared.HandleError(err, "Error while creating room")
	return room, err
}
func (rs roomSvc) AddMember(room *domain.Room, member *domain.Member) {
	client := db.GetClient()
	col := client.Database("chatroom").Collection("room")

	filter := bson.M{"roomId": room.RoomId}
	room.Members = append(room.Members, member)

	newRoom, err := col.ReplaceOne(context.TODO(), filter, &room)
	fmt.Println(newRoom)
	shared.HandleError(err, "Error while adding member to the room")
}

func (rs roomSvc) RemoveMember(room *domain.Room, username string) {
	// room.Lock()
	// defer room.Unlock()

	for i, member := range room.Members {
		if member.Username == username {

			room.Members = shared.RemoveIndex(room.Members, i)

			// update db
			client := db.GetClient()
			col := client.Database("chatroom").Collection("room")
			_, err := col.UpdateByID(context.TODO(), room.RoomId, room)
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
