package model

import (
	"chat-app/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"
)

type responseModel struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type messageBody struct {
	Message string `json:"message"`
}

func NewMessageBody(msg string) *messageBody {
	return &messageBody{Message: msg}
}

type roomDataBody struct {
	RoomId string   `json:"roomId"`
	Member []string `json:"member"`
}

func NewRoomDataBody(roomId string, members []string) *roomDataBody {
	return &roomDataBody{RoomId: roomId, Member: members}
}

type newRoomBody struct {
	RoomId string `json:"roomId"`
}

func NewNewRoomBody(roomId string) *newRoomBody {
	return &newRoomBody{RoomId: roomId}
}

type newMemberBody struct {
	Username string `json:"username"`
}

func NewNewMemberBody(username string) *newMemberBody {
	return &newMemberBody{Username: username}
}

type chatMessageSend struct {
	MessageType string `json:"messageType"`
	Sender      string `json:"sender"`
	Message     string `json:"message"`
	ChatId      string `json:"chatId"`
}

func NewChatMessageSend(messageType, sender, message, chatId string) *chatMessageSend {
	return &chatMessageSend{MessageType: messageType, Sender: sender, Message: message, ChatId: chatId}
}

type receiver struct {
	Channel string `json:"channel"`
	Uid     string `json:"uid"`
}

type ChatMessageReceive struct {
	Message string   `json:"message"`
	SendTo  receiver `json:"sendTo"`
}

func (cmr *ChatMessageReceive) UnmarshalJSON(data []byte) error {
	type alias ChatMessageReceive // alias is important else it will go in inf loop
	var chatMsg alias

	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return err
	}

	if chatMsg.SendTo.Channel != utils.CHANNEL_ROOM && chatMsg.SendTo.Channel != utils.CHANNEL_INDIVIDUAL {
		return errors.New("channel must be 'room' or 'individual'")
	}

	*cmr = ChatMessageReceive(chatMsg)
	return nil
}

// responses
func NewResponse(status int, data any, message string, success bool) responseModel {
	return responseModel{
		Status:  status,
		Data:    data,
		Message: message,
		Success: success,
	}
}

func StatusOK(data any) responseModel {
	return responseModel{
		Status:  http.StatusOK,
		Data:    data,
		Message: "OK",
		Success: true,
	}
}

func StatusInternalServerError() responseModel {
	return responseModel{
		Status:  http.StatusInternalServerError,
		Message: "Unhandled error occurred. Please try again later",
		Success: false,
	}
}

func StatusBadRequest(message string) responseModel {
	return responseModel{
		Status:  http.StatusBadRequest,
		Message: message,
		Success: false,
	}
}
