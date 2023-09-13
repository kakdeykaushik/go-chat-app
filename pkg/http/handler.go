package http

import (
	"chat-app/pkg/app"
	"chat-app/pkg/db"
	"chat-app/pkg/utils"
	"fmt"
	"log"
	"os"

	"net/http"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

// spins up server
func (s *server) Start() {
	port := getPort()
	fmt.Println("Server is running...")
	log.Fatalln(http.ListenAndServe(port, s))
}

// ServeHTTP to implement http.Handler interface
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	chatApp := getChatApp()

	// todo: add HTTP methods as well, use r.Method()
	// todo: add more feature APIs like create member etc. Think about it

	switch r.URL.Path {

	case "/", "/home":
		chatApp.Home(w, r)
	case "/room/new":
		chatApp.NewRoom(w, r)
	case "/room/join":
		chatApp.JoinRoom(w, r)
	case "/room/view":
		chatApp.ViewRoom(w, r)
	case "/room/leave":
		chatApp.LeaveRoom(w, r)
	case "/room/":
		chatApp.ChatRoom(w, r)
	default:
		http.NotFound(w, r)
	}
}

// helper
func getChatApp() *app.ChatApp {
	liveMemberConn := db.NewDB(utils.STORE_MEMORY)

	// use db.NewDB to get mongo store? - probably NO bcz here i dont want to provide entity T
	dbClient, err := db.GetClient()
	utils.FatalError(err, "error while connecting to DB")

	memberDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_MEMBER, "username")
	roomDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_ROOM, "roomId")

	memberSvc := app.NewMemberSvc(dbClient, memberDBConfig)
	roomSvc := app.NewRoomSvc(dbClient, roomDBConfig)

	return app.NewChatApp(liveMemberConn, memberSvc, roomSvc)

}

// helper
func getPort() string {
	port := os.Getenv("port")
	if port == "" {
		port = ":8000"
	}
	return port
}
