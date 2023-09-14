package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"chat-app/pkg/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type roomSvc struct {
	db     *mongo.Client
	config *db.Configuration
}

func NewRoomSvc(client *mongo.Client, config *db.Configuration) roomSvc {
	return roomSvc{db: client, config: config}
}

func (rs *roomSvc) GetRoom(roomId string) (*model.Room, error) {
	repo := db.NewMongoStore[entity.Room](rs.db, rs.config)
	roomEntity, err := repo.Get(roomId)
	if err != nil {
		return nil, err
	}
	room := model.EntityToModelRoom(roomEntity)
	return room, nil
}

func (rs *roomSvc) CreateRoom() (*model.Room, error) {
	roomId := generateRoomID()
	room := model.NewRoom(roomId, []*model.Member{})

	repo := db.NewMongoStore[entity.Room](rs.db, rs.config)
	roomEntity := model.ModelToEntityRoom(room)
	err := repo.Save(roomEntity)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *roomSvc) AddMember(room *model.Room, member *model.Member) error {

	if !rs.IsNewMember(room, member.Username) {
		return nil
	}

	room.Members = append(room.Members, member)

	repo := db.NewMongoStore[entity.Room](rs.db, rs.config)
	roomEntity := model.ModelToEntityRoom(room)
	err := repo.Update(room.RoomId, roomEntity)
	return err
}

func (rs *roomSvc) RemoveMember(room *model.Room, username string) error {
	for i, member := range room.Members {
		if member.Username == username {
			room.Members = utils.RemoveIndex(room.Members, i)
			// update db
			roomEntity := model.ModelToEntityRoom(room)
			repo := db.NewMongoStore[entity.Room](rs.db, rs.config)
			err := repo.Update(room.RoomId, roomEntity)
			return err
		}
	}

	return nil
}

func (rs *roomSvc) IsNewMember(room *model.Room, username string) bool {
	for _, member := range room.Members {
		if member.Username == username {
			return false
		}
	}
	return true
}

// helper
func generateRoomID() string {
	return uuid.New().String()[:5]
}
