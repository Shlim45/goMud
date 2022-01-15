package main

import (
	"fmt"
	"log"
	"net"
)

type Session struct {
	id   string
	conn net.Conn
}

func (s *Session) SessionId() string {
	return s.id
}

// TODO(jon): non-blocking write to session
func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

type SessionEvent struct {
	Session *Session
	Event   interface{}
}

type SessionCreatedEvent struct{}

type SessionDisconnectedEvent struct{}

type SessionInputEvent struct {
	input string
}

type SessionHandler struct {
	world        *World
	eventChannel <-chan SessionEvent
	users        map[string]*User
}

func NewSessionHandler(world *World, eventChannel <-chan SessionEvent) *SessionHandler {
	return &SessionHandler{
		world:        world,
		eventChannel: eventChannel,
		users:        map[string]*User{},
	}
}

func (s *Session) DisconnectSession(world *World, users map[string]*User) {
	sid := s.SessionId()
	user := users[sid]
	if user != nil {
		world.RemoveFromWorld(user.Character)

		delete(users, sid)

		if addr, ok := s.conn.RemoteAddr().(*net.TCPAddr); ok {
			ip := addr.IP.String()
			log.Println(fmt.Sprintf("Connection disconnected from '%s'", ip))
		}
	}
}

func (h *SessionHandler) Start() {
	for sessionEvent := range h.eventChannel {
		session := sessionEvent.Session
		sid := session.SessionId()

		switch event := sessionEvent.Event.(type) {

		case *SessionCreatedEvent:
			// create user

			character := &MOB{
				name:     generateName(),
				tickType: TICK_STOP,
			}

			user := &User{session, character, true}
			character.User = user

			character.Init()

			h.users[sid] = user
			h.world.HandleCharacterJoined(character)

			// TODO(jon): log user in, get account info etc

		case *SessionDisconnectedEvent:
			// remove user
			session.DisconnectSession(h.world, h.users)

		case *SessionInputEvent:

			user := h.users[sid]
			h.world.HandlePlayerInput(user.Character, event.input)
		}
	}
}
