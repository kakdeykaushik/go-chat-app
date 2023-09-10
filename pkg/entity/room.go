package entity

type Room struct {
	RoomId  string    `json:"roomId" bson:"roomId"`
	Members []*Member `json:"members" bson:"members"`
}
