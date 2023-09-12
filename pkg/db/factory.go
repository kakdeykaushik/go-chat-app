package db

import (
	"chat-app/pkg/types"
	"chat-app/pkg/utils"
)

// factory
func NewDB(name string) types.Storage {
	switch name {

	case utils.STORE_MEMORY:
		return newInMemoryStore()

	default:
		panic("invalid DB name") // panic? seriously?
	}
}
