package app

import (
	model "chat-app/pkg/models"
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

type ChatApp struct {
	memberService memberSvc
	roomService   roomSvc
}

/*
this should have its own
  - member service
  - room service
*/
func NewChatApp(memberService memberSvc, roomService roomSvc) *ChatApp {
	return &ChatApp{memberService: memberService, roomService: roomService}
}

func (c *ChatApp) Home(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(utils.HOMEPAGE)
	if err != nil {
		fmt.Println(err)
		resp := model.StatusInternalServerError()

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
		data := model.NewMessageBody("room does not exist")

		resp := model.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	// remove from DB and update in mem
	c.roomService.RemoveMember(room, username)

	message := fmt.Sprintf("%v left the room", username)
	go c.sendMessage("admin", room, utils.MT_LEAVE, message)

	data := model.NewMessageBody("room left successfully")

	resp := model.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)

}

func (c *ChatApp) ViewRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")

	room, err := c.roomService.GetRoom(roomId)
	if err != nil {
		data := model.NewMessageBody("Unable to get room")

		resp := model.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	var members = []string{}
	for _, member := range room.Members {
		members = append(members, member.Username)
	}

	data := model.NewRoomDataBody(roomId, members)

	resp := model.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) NewRoom(w http.ResponseWriter, r *http.Request) {
	room, err := c.roomService.CreateRoom()

	if err != nil {
		fmt.Println("Error while creating new room", err)
		data := model.NewMessageBody("unable to create room. please try again later")
		resp := model.StatusOK(data)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	data := model.NewNewRoomBody(room.RoomId)

	resp := model.StatusOK(data)
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
	data := model.NewMessageBody("user added to the room")

	resp := model.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) NewMember(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	member, err := c.memberService.CreateMember(username)

	if err != nil {
		fmt.Println("Error while creating new member", err)
		data := model.NewMessageBody("unable to create member. please try again later")
		resp := model.StatusOK(data)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	data := model.NewNewMemberBody(member.Username)
	resp := model.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (c *ChatApp) ChatRoom(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	socket, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		fmt.Println("Error while upgrading protocol: ", err)
		resp := model.StatusBadRequest("Unable to connect")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	defer socket.Close()

	username := r.URL.Query().Get("email")

	fmt.Println("Username: ", username)

	_, err = c.memberService.GetMember(username)
	if err != nil {
		fmt.Println("member not found")
		return
	}

	c.memberService.AddConn(username, socket)

	var message model.ChatMessageReceive

	// Read messages
	for {
		err := socket.ReadJSON(&message)
		fmt.Println("message: ", message)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			err := c.memberService.DeleteConn(username)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if err != nil {
			fmt.Println("Failed to read message", err)
			continue
		}

		go c.verifyAndSend(username, message)
	}
}

// This will do some checks if message can be sent or not
func (c *ChatApp) verifyAndSend(username string, message model.ChatMessageReceive) {

	if message.SendTo.Channel == utils.CHANNEL_ROOM {
		roomId := message.SendTo.Uid
		myRoom, err := c.roomService.GetRoom(roomId)
		if err != nil {
			fmt.Println(err)
			return
		}
		if c.roomService.IsNewMember(myRoom, username) {
			// not allowed to send in this room
			return
		}

		c.sendMessage(username, myRoom, utils.MT_MESSAGE, message.Message)
		return
	}

	if message.SendTo.Channel == utils.CHANNEL_INDIVIDUAL {
		usernameReceiver := message.SendTo.Uid

		receiver, err := c.memberService.GetMember(usernameReceiver)
		if err != nil {
			return
		}

		if usernameReceiver == username {
			c.sendMessage(username, receiver, utils.MT_MESSAGE, message.Message)
			return
		}

		me, err := c.memberService.GetMember(username)
		if err != nil {
			return
		}

		c.sendMessage(username, me, utils.MT_MESSAGE, message.Message)
		c.sendMessage(username, receiver, utils.MT_MESSAGE, message.Message)
	}
}

func (c *ChatApp) sendMessage(sender string, rec any, messageType string, message string) {
	switch receiver := rec.(type) {

	case *model.Room:
		r, err := c.roomService.GetRoom(receiver.RoomId)
		if err != nil {
			return
		}
		for _, member := range r.Members {
			go c.send(messageType, sender, message, member, receiver.RoomId)
		}

	case *model.Member:
		go c.send(messageType, sender, message, receiver, receiver.Username)

	default:
		fmt.Println("invalid receiver")
	}
}

func (c *ChatApp) send(messageType string, sender string, message string, receiver *model.Member, chatId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	sock, err := c.memberService.GetConn(receiver.Username)
	if err != nil {
		fmt.Printf("cannt send to %v - %v\n", receiver.Username, err)
		return
	}

	fmt.Printf("Sending message to: %v - %s\n", receiver.Username, message)
	msg := model.NewChatMessageSend(messageType, sender, message, chatId)
	err = sock.WriteJSON(msg)
	if err != nil {
		fmt.Println("Failed to Write message", err)
	}
}
