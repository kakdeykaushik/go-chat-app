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

// todo ; add NewMember fn
// todo ; add NewX fn for almost all struct
