package http

import (
	"chat-app/pkg/app"
	"chat-app/pkg/db"
	model "chat-app/pkg/models"
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

	liveMemberConn = db.NewDB(utils.STORE_MEMORY)
	roomToMember   = db.NewDB(utils.STORE_MEMORY)

	memberSvc = app.NewMemberSvc()
	roomSvc   = app.NewRoomSvc()
)

func Home(w http.ResponseWriter, r *http.Request) {
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

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")

	room, err := roomSvc.GetRoom(roomId)
	if err != nil {
		data := model.MessageBody{Message: "room does not exist"}

		resp := utils.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	roomSvc.RemoveMember(room, username)
	liveMemberConn.Delete(username)
	roomToMember.Save(roomId, room.Members)

	message := fmt.Sprintf("%v left the room", username)
	go sendMessage("admin", room, utils.MT_LEAVE, message)

	data := model.MessageBody{Message: "room left successfully"}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)

}

func ViewRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")

	room, err := roomSvc.GetRoom(roomId)
	if err != nil {
		var data = model.MessageBody{Message: "Unable to get room"}

		resp := utils.StatusOK(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)

		return
	}

	var members []string
	for _, member := range room.Members {
		members = append(members, member.Username)
	}

	var data = model.RoomDataBody{RoomId: roomId, Member: members}

	resp := utils.StatusOK(data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func NewRoom(w http.ResponseWriter, r *http.Request) {
	room, err := roomSvc.CreateRoom()

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

func ChatRoom(w http.ResponseWriter, r *http.Request) {
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
	myRoom, _ := roomSvc.GetRoom(roomId)

	if myRoom == nil {
		fmt.Println("room not found")

		data := model.MessageBody{Message: "room not found"}

		resp := utils.StatusOK(data)
		socket.WriteJSON(resp)

		return
	}

	// Check if current user is already in the room
	if roomSvc.IsNewMember(myRoom, username) {
		member := memberSvc.CreateMember(username, socket)
		roomSvc.AddMember(myRoom, member)

		// Update conn and inmem room data
		liveMemberConn.Save(username, socket)
		roomToMember.Save(roomId, myRoom.Members)

		go sendMessage(username, myRoom, utils.MT_NEWUSER, fmt.Sprintf("%s joined.", username))
	} else {
		// Update conn and inmem room data
		liveMemberConn.Save(username, socket)
		roomToMember.Save(roomId, myRoom.Members)
	}

	var message model.ChatMessageReceive

	// Read messages
	for {
		err := socket.ReadJSON(&message)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			// Remove member
			liveMemberConn.Delete(username)
			return
		}

		// Rare edge case
		_, err = liveMemberConn.Get(username)
		if err != nil {
			fmt.Println("Member not available")
			return
		}

		if err != nil {
			fmt.Println("Failed to read message", err)
		}
		go sendMessage(username, myRoom, utils.MT_MESSAGE, message.Message)
	}
}

func sendMessage(sender string, rec any, messageType string, message string) {
	switch receiver := rec.(type) {

	case *model.Room:
		rm, err := roomToMember.Get(receiver.RoomId)
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
			go send(messageType, sender, message, *member)
		}

	case *model.Member:
		go send(messageType, sender, message, *receiver)

	default:
		fmt.Println("invalid receiver")
	}
}

func send(messageType string, sender string, message string, receiver model.Member) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	msg := model.ChatMessageSend{MessageType: messageType, Sender: sender, Message: message}

	s, err := liveMemberConn.Get(receiver.Username)
	if err != nil {
		fmt.Printf("cannt send to %v\n", receiver.Username)
		return
	}

	sock, ok := s.(*websocket.Conn)
	if !ok {
		fmt.Println(s, ok)
		return
	}

	fmt.Printf("Sending message to: %v - %s\n", receiver.Username, message)
	err = sock.WriteJSON(msg)
	if err != nil {
		fmt.Println("Failed to Write message", err)
	}
}
