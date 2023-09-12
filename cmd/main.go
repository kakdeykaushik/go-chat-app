package main

import (
	ihttp "chat-app/pkg/http"
	"chat-app/pkg/utils"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	utils.FatalError(err, "error while loading .env")
}

func main() {
	server := ihttp.NewServer()
	server.Start()
}
