package db

import (
	"chat-app/pkg/domain"
	"chat-app/shared"
)

// factory
func NewDB(name string, configs ...any) domain.Storage {
	switch name {

	case shared.STORE_MEMORY:
		return newInMemoryStore()

	case shared.STORE_MONGO:
		var config = &Configuration{configs[0].(string), configs[1].(string)}
		return newMongoStore(config)

	default:
		panic("invalid DB name")

	}
}
