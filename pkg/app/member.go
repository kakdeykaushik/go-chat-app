package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"chat-app/pkg/utils"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberSvc struct {
	db     *mongo.Client
	config *db.Configuration
}

func NewMemberSvc(client *mongo.Client, config *db.Configuration) memberSvc {
	return memberSvc{db: client, config: config}
}

func (ms *memberSvc) CreateMember(username string, socket *websocket.Conn) (*model.Member, error) {

	m, err := ms.GetMember(username)
	// error - member does not exists
	if err == mongo.ErrNoDocuments {

		conn := model.Connection{Socket: socket, Mutex: &sync.Mutex{}}
		member := &model.Member{Username: username, Conn: &conn}

		memberEntity := utils.ModelToEntityMember(member)

		repo := db.NewMongoStore[entity.Member](ms.db, ms.config)
		err := repo.Save(memberEntity)

		if err != nil {
			return nil, err
		}
		return member, err
	}

	// error - some other error
	if err != nil {
		return nil, err
	}

	// everything ok
	return m, nil
}

func (ms *memberSvc) GetMember(username string) (*model.Member, error) {
	repo := db.NewMongoStore[entity.Member](ms.db, ms.config)
	memberEntity, err := repo.Get(username)

	if err != nil {
		return nil, err
	}

	member := utils.EntityToModelMember(memberEntity)

	return member, nil
}
