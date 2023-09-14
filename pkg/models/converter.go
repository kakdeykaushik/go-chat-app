package model

import (
	"chat-app/pkg/entity"
	"sync"
)

func ModelToEntityMember(m *Member) *entity.Member {
	var e = &entity.Member{}

	e.Username = m.Username
	return e
}

func EntityToModelMember(e *entity.Member) *Member {
	var m = &Member{}

	m.Username = e.Username
	m.Conn = &Connection{}
	m.Conn.Mutex = &sync.Mutex{}
	m.Conn.Socket = nil

	return m
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

	var m = &Room{}
	m.RoomId = e.RoomId
	m.Members = members

	return m

}
