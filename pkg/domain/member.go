package domain

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket *websocket.Conn
	sync.Mutex
}
type Member struct {
	Username string
	Conn     *Connection
}

// redundant
func (m *Member) UpdateConn(s *websocket.Conn) {}

// type MemberDB interface {
// 	Get(id string) (*Member, error)
// 	List() ([]*Member, error)
// 	Save(p *Member) error
// 	Delete(id string) error
// }
