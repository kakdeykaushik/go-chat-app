package http

import (
	"chat-app/pkg/app"
	"chat-app/pkg/db"
	"chat-app/pkg/domain"
	"chat-app/shared"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}

	store1 = db.NewDB(shared.STORE_MONGO, shared.DB_CHATROOM, shared.COLLECTION_MEMBER)
	store2 = db.NewDB(shared.STORE_MONGO, shared.DB_CHATROOM, shared.COLLECTION_ROOM)

	memberSvc = app.NewMemberSvc(store1)
	roomSvc   = app.NewRoomSvc(store2)

	liveConn = new(sync.Map)
	rooms    = new(sync.Map)
)

func Home(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("chat.html")
	shared.HandleError(err, "Failed to read HTML file")
	w.Write(content)
}

func NewRoom(w http.ResponseWriter, r *http.Request) {
	room, err := roomSvc.CreateRoom()
	shared.HandleError(err, "Error while creating new room")
	fmt.Printf("New room created: %s\n", room.RoomId)

	// todo; move struct outside
	var ss = struct {
		RoomId string `json:"roomId"`
	}{room.RoomId}

	json.NewEncoder(w).Encode(&ss)
}

func ChatRoom(w http.ResponseWriter, r *http.Request) {
	// debate - should be Upgraded right away or after some basic check like valid roomId
	socket, err := upgrader.Upgrade(w, r, w.Header())
	shared.HandleError(err, "Failed to Upgrade")
	defer socket.Close()

	fmt.Println("Protocol/IP: ", socket.LocalAddr().Network(), socket.LocalAddr().String())

	username := r.URL.Query().Get("email")
	roomId := r.URL.Query().Get("roomId")

	fmt.Printf("Username, roomId: %v, %v\n", username, roomId)

	myRoom, err := roomSvc.GetRoom(roomId)
	shared.HandleError(err, fmt.Sprintf("Error getting room with roomId %s", roomId))

	if myRoom == nil {
		fmt.Println("room not found")
		// todo: move struct outside
		msg := &struct {
			MessageType string `json:"messageType"`
			Sender      string `json:"sender"`
			Message     string `json:"message"`
		}{"information", "admin", "room not found"}
		socket.WriteJSON(msg)
		return
	}

	if roomSvc.IsNewMember(myRoom, username) {
		member := memberSvc.CreateMember(username, socket)
		// myRoom.Members = append(myRoom.Members, member)
		roomSvc.AddMember(myRoom, member)
		liveConn.Store(username, member.Conn.Socket)
		rooms.Store(roomId, myRoom.Members)

		go sendMessage(socket, myRoom, "new user", username, fmt.Sprintf("%s joined.", username))
	} else {
		member, err := memberSvc.GetMember(username)
		shared.HandleError(err, fmt.Sprintf("Error getting member with username %s. Creating new member", username))

		member.UpdateConn(socket) // redundant
		liveConn.Store(username, member.Conn.Socket)
		// rooms.Store(roomId, myRoom.Members)

		go sendMessage(socket, member, "already joined", username, fmt.Sprintf("%s is already connected.", username))
		// currentConnector := member

		// currentConnector := memberSvc.CreateMember(username, socket)
		// go sendMessage(socket, currentConnector, "already joined", username, fmt.Sprintf("%s is already connected.", username))
		// go sendMessage(socket, member, "new connection", username, fmt.Sprintf("%s is trying to connect from somewhere else", username))
	}
	fmt.Println(myRoom.Members)

	// todo: move struct
	var data struct {
		Message string // `json:"message"`
		// SentAt  time.Time
	}

	for {
		err := socket.ReadJSON(&data)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			roomSvc.RemoveMember(myRoom, username)

			go sendMessage(socket, myRoom, "leaving", username, fmt.Sprintf("%s left the room", username))
			return
		}

		shared.HandleError(err, "Failed to read message")
		go sendMessage(socket, myRoom, "message", username, data.Message)
	}
}

func sendMessage(socket *websocket.Conn, myRoom any, messageType string, sender string, message string) {

	fmt.Printf("room - %+v, sender %+v \n", myRoom, sender)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	switch receiver := myRoom.(type) {

	case *domain.Room:
		rm, _ := rooms.Load(receiver.RoomId)
		r := rm.([]*domain.Member)
		for _, member := range r {
			// member./
			send(messageType, sender, message, *member)
		}

		// for _, member := range receiver.Members {
		// 	// todo: move this struct outside
		// 	msg := &struct {
		// 		MessageType string `json:"messageType"`
		// 		Sender      string `json:"sender"`
		// 		Message     string `json:"message"`
		// 	}{messageType, sender, message}

		// 	s, ok := liveConn.Load(member.Username)
		// 	if !ok {
		// 		fmt.Println(s, ok)
		// 	}
		// 	sock, ok := s.(*websocket.Conn)
		// 	if !ok {
		// 		fmt.Println(s, ok)
		// 	}

		// 	fmt.Printf("Sending message to: %v - %s\n", member.Username, message)
		// 	err := sock.WriteJSON(msg)
		// 	// err := member.Conn.Socket.WriteJSON(msg)
		// 	shared.HandleError(err, "Failed to Write message")
		// }

	case *domain.Member:
		send(messageType, sender, message, *receiver)

	default:
		fmt.Println("invalid receiver")

	}

}

func send(messageType string, sender string, message string, receiver domain.Member) {
	// todo: move struct
	msg := &struct {
		MessageType string `json:"messageType"`
		Sender      string `json:"sender"`
		Message     string `json:"message"`
	}{messageType, sender, message}

	s, ok := liveConn.Load(receiver.Username)
	if !ok {
		fmt.Println(s, ok)
	}
	sock, ok := s.(*websocket.Conn)
	if !ok {
		fmt.Println(s, ok)
	}

	fmt.Printf("Sending message to: %v - %s\n", receiver.Username, message)
	err := sock.WriteJSON(msg)
	shared.HandleError(err, "Failed to Write message")
}
