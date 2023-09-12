package http

import (
	"chat-app/pkg/app"
	"chat-app/pkg/db"
	"chat-app/pkg/utils"
	"fmt"
	"log"

	"net/http"
)

//  todo - this is actuall handler - rename to handler.go

type server struct{}

func NewServer() *server {
	return &server{}
}

// spins up server
func (s *server) Start() {
	fmt.Println("Server is running...")
	log.Fatalln(http.ListenAndServe(":8000", s))
}

// ServeHTTP to implement http.Handler interface
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	chatApp := getChatApp()

	switch r.URL.Path {

	case "/":
		chatApp.Home(w, r)
	case "/room/new":
		chatApp.NewRoom(w, r)
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
	roomToMember := db.NewDB(utils.STORE_MEMORY)

	// use db.NewDB to get mongo store? - probably NO bcz here i dont want to provide entity T
	dbClient, err := db.GetClient()
	utils.FatalError(err, "error while connecting to DB")

	memberDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_MEMBER, "username")
	roomDBConfig := db.NewConfiguration(utils.DB_CHATROOM, utils.COLLECTION_ROOM, "roomId")

	memberSvc := app.NewMemberSvc(dbClient, memberDBConfig)
	roomSvc := app.NewRoomSvc(dbClient, roomDBConfig)

	return app.NewChatApp(liveMemberConn, roomToMember, memberSvc, roomSvc)

}
