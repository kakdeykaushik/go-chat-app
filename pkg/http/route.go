package http

import (
	"fmt"
	"log"

	"net/http"
)

func Start() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/room/new", NewRoom)
	http.HandleFunc("/room/view", ViewRoom)
	http.HandleFunc("/room/leave", LeaveRoom)
	http.HandleFunc("/room/", ChatRoom)

	fmt.Println("Server is running...")
	log.Fatalln(http.ListenAndServe(":8000", nil))

}
