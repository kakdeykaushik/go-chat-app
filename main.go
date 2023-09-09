package main

import (
	ihttp "chat-app/pkg/http"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", ihttp.Home)
	http.HandleFunc("/room/new", ihttp.NewRoom)
	http.HandleFunc("/room/view", ihttp.ViewRoom)
	http.HandleFunc("/room/", ihttp.ChatRoom)

	fmt.Println("Server is running...")
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
