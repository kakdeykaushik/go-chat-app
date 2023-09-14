package app

import (
	"chat-app/pkg/db"
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"chat-app/pkg/types"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberSvc struct {
	db       *mongo.Client
	config   *db.Configuration
	liveConn types.Storage
}

func NewMemberSvc(client *mongo.Client, config *db.Configuration, live types.Storage) memberSvc {
	return memberSvc{db: client, config: config, liveConn: live}
}

func (ms *memberSvc) CreateMember(username string, socket *websocket.Conn) (*model.Member, error) {

	m, err := ms.GetMember(username)
	// error - member does not exists
	if err == mongo.ErrNoDocuments {

		conn := model.Connection{Socket: socket, Mutex: &sync.Mutex{}}
		member := &model.Member{Username: username, Conn: &conn}

		memberEntity := model.ModelToEntityMember(member)

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

	member := model.EntityToModelMember(memberEntity)

	return member, nil
}

func (m *memberSvc) GetActive(username string) (*websocket.Conn, error) {
	v, err := m.liveConn.Get(username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	conn, ok := v.(*websocket.Conn)
	if !ok {
		fmt.Println("unable to convert")
	}

	return conn, nil
}

func (m *memberSvc) DeleteConn(username string) error {
	err := m.liveConn.Delete(username)
	return err
}

func (m *memberSvc) GetConn(username string) (*websocket.Conn, error) {

	v, err := m.liveConn.Get(username)

	if err != nil {
		fmt.Println("not found", err)
		return nil, err
	}

	conn, ok := v.(*websocket.Conn)

	if !ok {
		return nil, errors.New("unable to get connection")
	}

	return conn, nil
}

func (m *memberSvc) AddConn(username string, conn *websocket.Conn) error {
	err := m.liveConn.Save(username, conn)
	return err
}
