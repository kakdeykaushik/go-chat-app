package main

import (
	ihttp "chat-app/pkg/http"
	"fmt"
	"log"
	"net/http"
)

/*
var (
	store = db.NewDB(shared.STORE_MONGO)
	room  = app.NewRoomSvc(store)
)

func init() {
	room, err := room.CreateRoom()
	shared.HandleError(err, fmt.Sprintf("Error while creating room %s", "R1"))
	fmt.Printf("Room id: %s\n", room.RoomId)
}
*/

func main() {
	http.HandleFunc("/", ihttp.Home)
	http.HandleFunc("/room/new", ihttp.NewRoom)
	http.HandleFunc("/room/", ihttp.ChatRoom)

	fmt.Println("Server is running...")
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
