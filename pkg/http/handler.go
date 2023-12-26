package http

import (
	"chat-app/pkg/app"
	"chat-app/pkg/db"
	"chat-app/pkg/utils"
	"context"
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

	// actual routes

	switch {
	case r.URL.Path == "/healthcheck":
		w.Write([]byte("ok"))
	case r.URL.Path == "/", r.URL.Path == "/home" && r.Method == http.MethodGet:
		chatApp.Home(w, r)
	case r.URL.Path == "/room/new" && r.Method == http.MethodPost:
		chatApp.NewRoom(w, r)
	case r.URL.Path == "/room/join" && r.Method == http.MethodPost:
		chatApp.JoinRoom(w, r)
	case r.URL.Path == "/room/view" && r.Method == http.MethodGet:
		chatApp.ViewRoom(w, r)
	case r.URL.Path == "/room/leave" && r.Method == http.MethodDelete:
		chatApp.LeaveRoom(w, r)
	case r.URL.Path == "/room/":
		chatApp.ChatRoom(w, r)
	case r.URL.Path == "/member/add" && r.Method == http.MethodPost:
		chatApp.NewMember(w, r)
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

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		utils.FatalError(err, "error while pinging DB")
	}

	memberDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_MEMBER, "username")
	roomDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_ROOM, "roomId")

	memberSvc := app.NewMemberSvc(dbClient, memberDBConfig, liveMemberConn)
	roomSvc := app.NewRoomSvc(dbClient, roomDBConfig)

	return app.NewChatApp(memberSvc, roomSvc)

}

// helper
func getPort() string {
	port := os.Getenv("port")
	if port == "" {
		port = ":8000"
	}
	return port
}
