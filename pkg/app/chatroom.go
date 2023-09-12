package app

import (
	model "chat-app/pkg/models"
	"chat-app/pkg/types"
	"chat-app/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")
		return !(username == "" || roomId == "")
	}}
)

type ChatApp struct {
	liveMemberConn types.Storage
	roomToMember   types.Storage
	memberService  memberSvc
	roomService    roomSvc
}

/*
this should have its own
- liveMemberConn
- roomToMember
- and rest as usual
  - db client
  - configs (for below services)
  - member service
  - room service
*/
func NewChatApp(liveMemberConn types.Storage, roomToMember types.Storage, memberService memberSvc, roomService roomSvc) *ChatApp {
	return &ChatApp{liveMemberConn: liveMemberConn, roomToMember: roomToMember, memberService: memberService, roomService: roomService}
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
	c.roomToMember.Save(roomId, room.Members)

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

	username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")
	fmt.Printf("Username, roomId: %v, %v\n", username, roomId)

	// Get room
	myRoom, err := c.roomService.GetRoom(roomId)

	if err != nil {
		fmt.Println("room not found", err)
		data := model.MessageBody{Message: "room not found"}

		resp := utils.StatusOK(data)
		socket.WriteJSON(resp)

		return
	}

	// Check if current user is already in the room
	if c.roomService.IsNewMember(myRoom, username) {
		member, err := c.memberService.CreateMember(username, socket)
		if err != nil {
			fmt.Println("error while creating member", err)
			return // internal server error
		}

		err = c.roomService.AddMember(myRoom, member)
		if err != nil {
			fmt.Println("error while adding member", err)
			return // internal server error
		}

		// Update conn and inmem room data
		c.liveMemberConn.Save(username, socket)
		c.roomToMember.Save(roomId, myRoom.Members)

		go c.sendMessage(username, myRoom, utils.MT_NEWUSER, fmt.Sprintf("%s joined.", username))
	} else {
		// Update conn and inmem room data
		c.liveMemberConn.Save(username, socket)
		c.roomToMember.Save(roomId, myRoom.Members)
	}

	var message model.ChatMessageReceive

	// Read messages
	for {
		err := socket.ReadJSON(&message)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			// Remove member
			c.liveMemberConn.Delete(username)
			return
		}

		// Rare edge case
		_, err = c.liveMemberConn.Get(username)
		if err != nil {
			fmt.Println("Member not available")
			return
		}

		if err != nil {
			fmt.Println("Failed to read message", err)
			continue
		}
		go c.sendMessage(username, myRoom, utils.MT_MESSAGE, message.Message)
	}
}

func (c *ChatApp) sendMessage(sender string, rec any, messageType string, message string) {
	switch receiver := rec.(type) {

	case *model.Room:
		rm, err := c.roomToMember.Get(receiver.RoomId)
		if err != nil {
			fmt.Println("Unable to get members")
			return
		}

		members, ok := rm.([]*model.Member)
		if !ok {
			fmt.Println("not parsed to []*model.Member")
			return
		}

		for _, member := range members {
			go c.send(messageType, sender, message, member)
		}

	case *model.Member:
		go c.send(messageType, sender, message, receiver)

	default:
		fmt.Println("invalid receiver")
	}
}

func (c *ChatApp) send(messageType string, sender string, message string, receiver *model.Member) {
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
	msg := model.ChatMessageSend{MessageType: messageType, Sender: sender, Message: message}
	err = sock.WriteJSON(msg)
	if err != nil {
		fmt.Println("Failed to Write message", err)
	}
}
