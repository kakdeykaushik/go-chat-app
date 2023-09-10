package db

import (
	"chat-app/pkg/types"
	"chat-app/pkg/utils"
)

// factory
func NewDB(name string, configs ...string) types.Storage {
	switch name {

	case utils.STORE_MEMORY:
		return newInMemoryStore()

	case utils.STORE_MONGO:
		var config = &Configuration{configs[0], configs[1]}
		return newMongoStore(config)

	default:
		panic("invalid DB name")
	}
}
