package main

import (
	"fmt"
	"log"
	"net"
)

type Session struct {
	id     string
	conn   net.Conn
	status SessionStatus
}

func (s *Session) SessionId() string {
	return s.id
}

func (s *Session) Status() SessionStatus {
	return s.status
}

func (s *Session) SetStatus(newStatus SessionStatus) {
	s.status = newStatus
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
	eventChannel <-chan SessionEvent
	users        map[string]*User
	library      *MudLib
}

func NewSessionHandler(library *MudLib, eventChannel <-chan SessionEvent) *SessionHandler {
	return &SessionHandler{
		library:      library,
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

type SessionStatus uint8

const (
	DEFAULT = iota
	USERNAME
	PASSWORD
	MENU
	SELECT
	CREATE
	INGAME
	QUIT
)

func (h *SessionHandler) Start() {
	for sessionEvent := range h.eventChannel {
		session := sessionEvent.Session
		sid := session.SessionId()

		switch event := sessionEvent.Event.(type) {

		case *SessionCreatedEvent:
			user := &User{session, nil, nil, true}

			h.users[sid] = user

			session.WriteLine(MudASCIILogo())
			session.WriteLine("Welcome to Darkness Falls.\r\n")
			session.WriteLine("\r\nUsername: ")
			session.SetStatus(USERNAME)

		case *SessionDisconnectedEvent:
			session.DisconnectSession(h.library.world, h.users)

		case *SessionInputEvent:
			user := h.users[sid]

			if session.Status() == INGAME {
				h.library.world.HandlePlayerInput(user.Character, event.input, h.library)
			} else {
				h.library.world.HandleUserLogin(user, event.input)
			}
		}
	}
}
