package model

import (
	"chat-app/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"
)

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

func (cmr *ChatMessageReceive) UnmarshalJSON(data []byte) error {
	var chatMsg ChatMessageReceive

	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return err
	}

	if chatMsg.SendTo.Channel != utils.CHANNEL_ROOM && chatMsg.SendTo.Channel != utils.CHANNEL_INDIVIDUAL {
		return errors.New("channel must be 'room' or 'individual'")
	}

	*cmr = ChatMessageReceive(chatMsg)
	return nil
}

func NewResponse(status int, data any, message string, success bool) ResponseModel {
	return ResponseModel{
		Status:  status,
		Data:    data,
		Message: message,
		Success: success,
	}
}

func StatusOK(data any) ResponseModel {
	return ResponseModel{
		Status:  http.StatusOK,
		Data:    data,
		Message: "OK",
		Success: true,
	}
}

func StatusInternalServerError() ResponseModel {
	return ResponseModel{
		Status:  http.StatusInternalServerError,
		Message: "Unhandled error occurred. Please try again later",
		Success: false,
	}
}

func StatusBadRequest(message string) ResponseModel {
	return ResponseModel{
		Status:  http.StatusBadRequest,
		Message: message,
		Success: false,
	}
}
