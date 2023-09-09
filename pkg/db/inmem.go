package db

// import (
// 	"chat-app/pkg/domain"
// 	"fmt"

// 	"github.com/pkg/errors"
// )

// type inMemoryStorage struct {
// 	items map[any]any
// }

// var singleInMemInstance *inMemoryStorage

// func newInMemoryStore() domain.MemberDB {
// 	if singleInMemInstance == nil {
// 		singleInMemInstance = &inMemoryStorage{items: make(map[any]any)}
// 	}

// 	return singleInMemInstance
// }

// func (s *inMemoryStorage) Get(K string) (*domain.Member, error) {
// 	fmt.Println(s.items)
// 	value, ok := s.items[K]
// 	if ok {
// 		return value.(*domain.Member), nil
// 	}
// 	return nil, errors.Errorf("no data found with key: %v", K)
// }

// func (s *inMemoryStorage) List() ([]*domain.Member, error) {
// 	var values []*domain.Member
// 	for _, item := range s.items {
// 		values = append(values, item.(*domain.Member))
// 	}

// 	return values, nil
// }

// func (s *inMemoryStorage) Save(m *domain.Member) error {
// 	s.items[m.Username] = m
// 	fmt.Println(s.items)
// 	return nil
// }

// func (s *inMemoryStorage) Delete(K string) error {
// 	delete(s.items, K)
// 	return nil
// }
