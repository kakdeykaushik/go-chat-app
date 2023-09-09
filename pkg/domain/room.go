package domain

type Room struct {
	RoomId  string    `json:"roomId" bson:"roomId"`
	Members []*Member `json:"members" bson:"members"`
	// sync.Mutex
}

// type RoomDB interface {
// 	Get(id string) (*Room, error)
// 	List() ([]*Room, error)
// 	Save(r *Room) error
// 	Delete(id string) error
// }
