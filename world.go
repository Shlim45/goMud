package main

import (
	"fmt"
	"log"
	"strings"
)

type World struct {
	accounts    map[string]*Account
	charClasses map[string]CharClass
	characters  map[string]*MOB
	areas       map[string]*Area
	rooms       []*Room
	db          DBConnection
}

func NewWorld(db DBConnection) *World {
	return &World{
		accounts:    make(map[string]*Account),
		charClasses: make(map[string]CharClass),
		characters:  make(map[string]*MOB),
		areas:       make(map[string]*Area),
		db:          db,
	}
}

func (w *World) AddArea(area *Area) {
	_, exists := w.areas[area.Name]
	if !exists {
		w.areas[area.Name] = area
	}
}

func (w *World) Init() {
	w.areas["Darkness Falls"] = &Area{
		Name:  "Darkness Falls",
		Realm: REALM_IMMORTAL,
	}

	w.rooms = []*Room{
		{
			Id:   "A",
			Desc: "standing in a room with a sign that has the letter A on it.",
			Area: w.areas["Darkness Falls"],
			Links: []*RoomLink{
				{
					Verb:   "east",
					RoomId: "B",
				},
			},
			Portals: []*RoomLink{},
		},
		{
			Id:   "B",
			Desc: "standing in a room with a sign that has the letter B on it.",
			Area: w.areas["Darkness Falls"],
			Links: []*RoomLink{
				{
					Verb:   "west",
					RoomId: "A",
				},
				{
					Verb:   "east",
					RoomId: "C",
				},
			},
			Portals: []*RoomLink{
				{
					Verb:   "gate",
					RoomId: "D",
				},
			},
		},
		{
			Id:   "C",
			Desc: "standing in a room with a sign that has the letter C on it.",
			Area: w.areas["Darkness Falls"],
			Links: []*RoomLink{
				{
					Verb:   "west",
					RoomId: "B",
				},
			},
			Portals: []*RoomLink{},
		},
		{
			Id:    "D",
			Desc:  "standing in a room hidden behind a gate.  There is a sign that has the letter D on it.",
			Area:  w.areas["Darkness Falls"],
			Links: []*RoomLink{},
			Portals: []*RoomLink{
				{
					Verb:   "gate",
					RoomId: "B",
				},
			},
		},
	}
	areaRooms := w.rooms
	w.areas["Darkness Falls"].Rooms = areaRooms
}

func (w *World) HandleCharacterJoined(character *MOB) {
	w.characters[character.Name()] = character
	w.rooms[0].AddMOB(character)

	character.SendMessage("Welcome to Darkness Falls\n\r", true)
	character.Room.ShowRoom(character)
	character.Room.ShowOthers(character, nil, fmt.Sprintf("%s appears in a puff of smoke.", character.Name()))

	log.Println(fmt.Sprintf("Player login: %s", character.Name()))
}

func (w *World) RemoveFromWorld(character *MOB) {
	room := character.Room
	room.RemoveMOB(character)
	room.Show(nil, fmt.Sprintf("%s disappears in a puff of smoke.", character.Name()))

	log.Println(fmt.Sprintf("Player logout: %s", character.Name()))
}

func (w *World) Broadcast(msg string) {
	for _, player := range w.characters {
		player.SendMessage(msg, true)
	}
}

func (w *World) GetRoomById(id string) *Room {
	for _, r := range w.rooms {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func (w *World) MoveCharacter(character *MOB, to *Room) {
	character.Room.RemoveMOB(character)
	to.AddMOB(character)
	to.ShowRoom(character)
	character.adjFat(-2, character.maxState().fat())
}

func (w *World) HandlePlayerInput(player *MOB, input string, library *MudLib) {
	tokens := strings.Split(input, " ")
	success := false
	cmd := library.FindCommand(tokens[0])

	for _, link := range player.Room.Links {
		if link.Verb == tokens[0] || strings.HasPrefix(link.Verb, tokens[0]) {
			if player.curState().fat() < 2 {
				player.SendMessage("You are too tired!", true)
				return
			}

			target := w.GetRoomById(link.RoomId)
			if target != nil {
				player.SendMessage(fmt.Sprintf("You travel %s.", link.Verb), true)
				player.Room.ShowOthers(player, nil, fmt.Sprintf("%s went %s.", player.Name(), link.Verb))
				w.MoveCharacter(player, target)
				player.Room.ShowOthers(player, nil, fmt.Sprintf("%s just came in.", player.Name()))
				success = true
				break
			}
		}
	}

	if !success {
		if cmd != nil {
			go cmd.ExecuteCmd(player, tokens, w, library)
		} else {
			player.SendMessage(fmt.Sprintf("Huh?  Command '%s' not found.", tokens[0]), true)
		}
	}
}
