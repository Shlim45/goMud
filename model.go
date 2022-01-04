package main

import (
	"fmt"
	"math/rand"
)

type SessionEvent struct {
	Session *Session
	Event   interface{}
}

type SessionCreatedEvent struct{}

type SessionDisconnectedEvent struct{}

type SessionInputEvent struct {
	input string
}

type Entity struct {
	entityId string
}

func (e *Entity) EntityId() string {
	return e.entityId
}

type User struct {
	Session   *Session
	Character *Character
}

type Character struct {
	Name string
	User *User
	Room *Room
}

func (c *Character) SendMessage(msg string) {
	c.User.Session.WriteLine(msg)
}

type RoomLink struct {
	Verb   string
	RoomId string
}

type Room struct {
	Id    string
	Desc  string
	Links []*RoomLink

	Characters []*Character
}

func (r *Room) AddCharacter(character *Character) {
	r.Characters = append(r.Characters, character)
	character.Room = r
}

func (r *Room) RemoveCharacter(character *Character) {
	character.Room = nil

	var characters []*Character
	for _, c := range r.Characters {
		if c != character {
			characters = append(characters, c)
		}
	}
	r.Characters = characters
}

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}

//
//type MessageEvent struct {
//	msg string
//}
//
//type MoveEvent struct {
//	dir string
//}
//
//type UserJoinedEvent struct {
//}
