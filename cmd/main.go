package main

import (
	ihttp "chat-app/pkg/http"
	"chat-app/pkg/utils"

	"github.com/joho/godotenv"
)

func init() {
	// loan env or crash
	err := godotenv.Load("../.env")
	utils.HandleError(err, "error while loading env")
}

func main() {
	ihttp.Start()
}
