package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket *websocket.Conn
	*sync.Mutex
}

type Member struct {
	Username string
	Conn     *Connection
}

// todo ;
// 1. add NewMember fn
// 2. add NewX fn for almost all struct
