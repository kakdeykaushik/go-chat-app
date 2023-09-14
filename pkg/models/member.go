package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {
	Socket *websocket.Conn
	*sync.Mutex
}

type Member struct {
	Username string
	Conn     *connection
}

func NewMember(username string, socket *websocket.Conn) *Member {
	conn := &connection{Socket: socket, Mutex: &sync.Mutex{}}
	return &Member{Username: username, Conn: conn}
}
