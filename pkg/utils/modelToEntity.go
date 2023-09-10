package utils

import (
	"chat-app/pkg/entity"
	model "chat-app/pkg/models"
	"sync"
)

func ModelToEntityMember(m *model.Member) *entity.Member {
	var e = &entity.Member{}

	e.Username = m.Username
	return e
}

func EntityToModelMember(e *entity.Member) *model.Member {
	var m = &model.Member{}

	m.Username = e.Username
	m.Conn = &model.Connection{}
	m.Conn.Mutex = &sync.Mutex{}
	m.Conn.Socket = nil

	return m
}

func ModelToEntityRoom(r *model.Room) *entity.Room {

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
func EntityToModelRoom(e *entity.Room) *model.Room {

	var members = []*model.Member{}
	for _, member := range e.Members {
		modelMember := EntityToModelMember(member)
		members = append(members, modelMember)
	}

	var m = &model.Room{}
	m.RoomId = e.RoomId
	m.Members = members

	return m

}
