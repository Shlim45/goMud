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
