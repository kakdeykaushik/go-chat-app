package entity

type Room struct {
	RoomId  string    `bson:"roomId"`
	Members []*Member `bson:"members"`
}
