package model

type Room struct {
	RoomId  string    `json:"roomId" bson:"roomId"`
	Members []*Member `json:"members" bson:"members"`
}

func NewRoom(roomId string, members []*Member) *Room {
	return &Room{RoomId: roomId, Members: members}
}
