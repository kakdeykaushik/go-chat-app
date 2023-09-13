package app

import (
	model "chat-app/pkg/models"
	"chat-app/pkg/types"
	"chat-app/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		username := r.URL.Query().Get("email")
		return strings.Contains(username, ".com")
	}}
)

// todo; move liveMemberConn to memberService and think how to reduce DB calls
type ChatApp struct {
	liveMemberConn types.Storage // stores K:V as username : ws 			 - same below this can moved to memberService
	memberService  memberSvc
	roomService    roomSvc
}

/*
this should have its own
- liveMemberConn
- and rest as usual
  - db client
  - configs (for below services)
  - member service
  - room service
*/
func NewChatApp(liveMemberConn types.Storage, memberService memberSvc, roomService roomSvc) *ChatApp {
	return &ChatApp{liveMemberConn: liveMemberConn, memberService: memberService, roomService: roomService}
}

func (c *ChatApp) Home(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(utils.HOMEPAGE)
	if err != nil {
		fmt.Println(err)
		resp := utils.StatusInternalServerError()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	w.Write(content)
}

func (c *ChatApp) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")

	room, err := c.roomService.GetRoom(roomId)
	if err != nil {
		data := model.MessageBody{Message: "room does not exist"}

		resp := utils.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	// remove from DB and update in mem
	c.roomService.RemoveMember(room, username)
	c.liveMemberConn.Delete(username)

	message := fmt.Sprintf("%v left the room", username)
	go c.sendMessage("admin", room, utils.MT_LEAVE, message)

	data := model.MessageBody{Message: "room left successfully"}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)

}

func (c *ChatApp) ViewRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")

	room, err := c.roomService.GetRoom(roomId)
	if err != nil {
		var data = model.MessageBody{Message: "Unable to get room"}

		resp := utils.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	var members = []string{}
	for _, member := range room.Members {
		members = append(members, member.Username)
	}

	var data = model.RoomDataBody{RoomId: roomId, Member: members}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) NewRoom(w http.ResponseWriter, r *http.Request) {
	room, err := c.roomService.CreateRoom()

	if err != nil {
		fmt.Println("Error while creating new room", err)
		data := model.MessageBody{Message: "unable to create room. please try again later"}
		resp := utils.StatusOK(data)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	data := model.NewRoomBody{RoomId: room.RoomId}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) JoinRoom(w http.ResponseWriter, r *http.Request) {
	username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")

	myRoom, err := c.roomService.GetRoom(roomId)
	if err != nil {
		fmt.Println(err)
		return
	}

	member, err := c.memberService.GetMember(username)
	if err != nil {
		return
	}
	err = c.roomService.AddMember(myRoom, member)
	if err != nil {
		return
	}
	data := &model.MessageBody{Message: "user added to the room"}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) ChatRoom(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	socket, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		fmt.Println("Error while upgrading protocol: ", err)
		resp := utils.StatusBadRequest("Unable to connect")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	defer socket.Close()

	username := r.URL.Query().Get("email")
	c.liveMemberConn.Save(username, socket)

	fmt.Println("Username: ", username)

	var message model.ChatMessageReceive

	// Read messages
	for {
		err := socket.ReadJSON(&message)

		fmt.Println("message: ", message)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			// Remove member
			c.liveMemberConn.Delete(username)
			return
		}

		if err != nil {
			fmt.Println("Failed to read message", err)
			continue
		}

		if message.SendTo.Channel == utils.CHANNEL_ROOM {
			roomId := message.SendTo.Uid
			myRoom, err := c.roomService.GetRoom(roomId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if c.roomService.IsNewMember(myRoom, username) {
				return
			}

			go c.sendMessage(username, myRoom, utils.MT_MESSAGE, message.Message)
		}

		if message.SendTo.Channel == utils.CHANNEL_INDIVIDUAL {
			usernameReceiver := message.SendTo.Uid

			receiver, err := c.memberService.GetMember(usernameReceiver)

			if err != nil {
				continue
			}
			me, err := c.memberService.GetMember(username)
			if err != nil {
				continue
			}

			go c.sendMessage(username, me, utils.MT_MESSAGE, message.Message)
			go c.sendMessage(username, receiver, utils.MT_MESSAGE, message.Message)
		}

	}
}

func (c *ChatApp) sendMessage(sender string, rec any, messageType string, message string) {
	switch receiver := rec.(type) {

	case *model.Room:
		r, _ := c.roomService.GetRoom(receiver.RoomId)
		for _, member := range r.Members {
			go c.send(messageType, sender, message, member, receiver.RoomId)
		}

	case *model.Member:
		go c.send(messageType, sender, message, receiver, receiver.Username)

	default:
		fmt.Println("invalid receiver")
	}
}

func (c *ChatApp) send(messageType string, sender string, message string, receiver *model.Member, roomId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	s, err := c.liveMemberConn.Get(receiver.Username)
	if err != nil {
		fmt.Printf("cannt send to %v\n", receiver.Username)
		return
	}

	sock, ok := s.(*websocket.Conn)
	if !ok {
		fmt.Println(s, ok, "oeidfh")
		return
	}

	fmt.Printf("Sending message to: %v - %s\n", receiver.Username, message)
	msg := model.ChatMessageSend{MessageType: messageType, Sender: sender, Message: message, RoomId: roomId}
	err = sock.WriteJSON(msg)
	if err != nil {
		fmt.Println("Failed to Write message", err)
	}
}

/*
to pivot app from "room" to "room(s) and DM(s)"
URL check for roomId should be dropped and some kind of Auth can be implemented

then - on .ReadJSON data should contain to whom message should be sent to room(roomId) or DM(username)
*/
