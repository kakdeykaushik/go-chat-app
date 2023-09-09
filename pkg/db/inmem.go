package db

import (
	"chat-app/pkg/domain"
	"sync"

	"github.com/pkg/errors"
)

type inMemoryStorage struct {
	items *sync.Map
}

var singleInMemInstance *inMemoryStorage

// singleton
func newInMemoryStore() domain.Storage {
	if singleInMemInstance == nil {
		singleInMemInstance = &inMemoryStorage{items: new(sync.Map)}
	}
	return singleInMemInstance
}

func (s *inMemoryStorage) Get(K string) (any, error) {
	value, ok := s.items.Load(K)
	if ok {
		return value, nil
	}
	return nil, errors.Errorf("no data found with key: %v", K)
}

func (s *inMemoryStorage) List() ([]any, error) {
	var values []any
	s.items.Range(func(key, value interface{}) bool {
		values = append(values, value)
		return true
	})
	return values, nil
}

func (s *inMemoryStorage) Save(K, V any) error {
	s.items.Store(K, V)
	return nil
}

func (s *inMemoryStorage) Delete(K string) error {
	s.items.Delete(K)
	return nil
}
