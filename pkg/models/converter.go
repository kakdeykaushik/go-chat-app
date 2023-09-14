package model

import (
	"chat-app/pkg/entity"
)

func ModelToEntityMember(m *Member) *entity.Member {
	var e = &entity.Member{}

	e.Username = m.Username
	return e
}

func EntityToModelMember(e *entity.Member) *Member {
	return NewMember(e.Username, nil)
}

func ModelToEntityRoom(r *Room) *entity.Room {

	var members = []*entity.Member{}
	for _, member := range r.Members {
		entityMember := ModelToEntityMember(member)
		members = append(members, entityMember)
	}

	var e = &entity.Room{}
	e.RoomId = r.RoomId
	e.Members = members
	return e

}
func EntityToModelRoom(e *entity.Room) *Room {

	var members = []*Member{}
	for _, member := range e.Members {
		modelMember := EntityToModelMember(member)
		members = append(members, modelMember)
	}

	return NewRoom(e.RoomId, members)

}
