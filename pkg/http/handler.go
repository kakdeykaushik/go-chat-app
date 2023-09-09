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

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")
		return !(username == "" || roomId == "")
	}}

	liveMemberConn = db.NewDB(shared.STORE_MEMORY)
	roomToMember   = db.NewDB(shared.STORE_MEMORY)

	memberSvc = app.NewMemberSvc()
	roomSvc   = app.NewRoomSvc()
)

func Home(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(shared.HOMEPAGE)
	shared.HandleError(err, "Failed to read HTML file")
	w.Write(content)
}

func ViewRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")

	room, err := roomSvc.GetRoom(roomId)
	if err != nil {
		w.Write([]byte("Something went wrong"))
		fmt.Printf("Error: %v", err)
		return
	}

	type Resp struct {
		RoomId string   `json:"roomId"`
		Member []string `json:"member"`
	}

	var members []string
	for _, member := range room.Members {
		members = append(members, member.Username)
	}

	var resp = Resp{RoomId: roomId, Member: members}
	json.NewEncoder(w).Encode(&resp)
}

func NewRoom(w http.ResponseWriter, r *http.Request) {
	room, err := roomSvc.CreateRoom()
	shared.HandleError(err, "Error while creating new room")

	// todo; move struct outside
	type Resp struct {
		RoomId string `json:"roomId"`
	}

	var resp = Resp{RoomId: room.RoomId}
	json.NewEncoder(w).Encode(&resp)
}

func ChatRoom(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	socket, err := upgrader.Upgrade(w, r, w.Header())
	shared.HandleError(err, "Failed to Upgrade")
	defer socket.Close()

	username, roomId := r.URL.Query().Get("email"), r.URL.Query().Get("roomId")
	fmt.Printf("Username, roomId: %v, %v\n", username, roomId)

	// Get room
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
		_ = socket.WriteJSON(msg)
		return
	}

	// Check if current user is already in the room
	if roomSvc.IsNewMember(myRoom, username) {
		member := memberSvc.CreateMember(username, socket)
		roomSvc.AddMember(myRoom, member)

		// Update conn and inmem room data
		liveMemberConn.Save(username, member.Conn.Socket)
		roomToMember.Save(roomId, myRoom.Members)

		go sendMessage(username, myRoom, shared.MT_NEWUSER, fmt.Sprintf("%s joined.", username))
	} else {
		// Update conn and inmem room data
		liveMemberConn.Save(username, socket)
		roomToMember.Save(roomId, myRoom.Members)
	}

	// todo: move struct
	var data struct {
		Message string // `json:"message"`
		// SentAt  time.Time
	}

	// Read messages
	for {
		err := socket.ReadJSON(&data)

		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			// Remove member
			liveMemberConn.Delete(username)
			roomSvc.RemoveMember(myRoom, username)
			roomToMember.Save(roomId, myRoom.Members)

			// waiting not required for this goroutine as this is fire and forget task. think about it again.
			// this message is not required as per biz logic since this is going offline. Instead create feature for member to exit room and then show this msg
			// go sendMessage(username, myRoom, shared.MT_LEAVE, fmt.Sprintf("%s left the room", username))
			return
		}

		shared.HandleError(err, "Failed to read message")
		go sendMessage(username, myRoom, shared.MT_MESSAGE, data.Message)
	}
}

func sendMessage(sender string, rec any, messageType string, message string) {
	switch receiver := rec.(type) {

	case *domain.Room:
		rm, _ := roomToMember.Get(receiver.RoomId)
		members, ok := rm.([]*domain.Member)
		if !ok {
			fmt.Println("not parsed to []*domain.Member")
		}
		for _, member := range members {
			go send(messageType, sender, message, *member)
		}

	case *domain.Member:
		go send(messageType, sender, message, *receiver)

	default:
		fmt.Println("invalid receiver")
	}
}

func send(messageType string, sender string, message string, receiver domain.Member) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	// todo: move struct
	msg := &struct {
		MessageType string `json:"messageType"`
		Sender      string `json:"sender"`
		Message     string `json:"message"`
	}{messageType, sender, message}

	s, err := liveMemberConn.Get(receiver.Username)
	if err != nil {
		fmt.Printf("cannt send to %v\n", receiver.Username)
		return
	}

	sock, ok := s.(*websocket.Conn)
	if !ok {
		fmt.Println(s, ok)
	}

	fmt.Printf("Sending message to: %v - %s\n", receiver.Username, message)
	err = sock.WriteJSON(msg)
	shared.HandleError(err, "Failed to Write message")
}
