package db

import (
	"chat-app/pkg/domain"
	"chat-app/shared"
)

// factory
func NewDB(name string, configs ...string) domain.Storage {
	switch name {

	case shared.STORE_MEMORY:
		return newInMemoryStore()

	case shared.STORE_MONGO:
		var config = &Configuration{configs[0], configs[1]}
		return newMongoStore(config)

	default:
		panic("invalid DB name")
	}
}
