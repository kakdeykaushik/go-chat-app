package model

type ResponseModel struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type MessageBody struct {
	Message string `json:"message"`
}

type RoomDataBody struct {
	RoomId string   `json:"roomId"`
	Member []string `json:"member"`
}

type NewRoomBody struct {
	RoomId string `json:"roomId"`
}

type ChatMessageSend struct {
	MessageType string `json:"messageType"`
	Sender      string `json:"sender"`
	Message     string `json:"message"`
	RoomId      string `json:"roomId"`
}

type Receiver struct {
	Channel string `json:"channel"`
	Uid     string `json:"uid"`
}

type ChatMessageReceive struct {
	Message string   `json:"message"`
	SendTo  Receiver `json:"sendTo"`
	// SentAt  time.Time // todo
}

// todo - implement "UnmarshalJSON" to restrict Channel to 2 types
