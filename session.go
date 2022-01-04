package main

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

func (h *SessionHandler) Start() {
	for sessionEvent := range h.eventChannel {
		session := sessionEvent.Session
		sid := session.SessionId()

		switch event := sessionEvent.Event.(type) {

		case *SessionCreatedEvent:
			// create user

			character := &Character{
				Name: generateName(),
			}
			user := &User{session, character}
			character.User = user

			h.users[sid] = user
			h.world.HandleCharacterJoined(character)

			// TODO(jon): log user in, get account info etc

		case *SessionDisconnectedEvent:
			// remove user
		case *SessionInputEvent:

			user := h.users[sid]
			h.world.HandleCharacterInput(user.Character, event.input)
		}
	}
}
