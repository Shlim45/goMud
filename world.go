package main

import (
	"fmt"
	"log"
	"strings"
)

type World struct {
	characters []*Character
	rooms      []*Room
}

func NewWorld() *World {
	return &World{}
}

func (w *World) Init() {
	w.rooms = []*Room{
		{
			Id:   "A",
			Desc: "You're standing in a room with a sign that has the letter A on it.",
			Links: []*RoomLink{
				{
					Verb:   "east",
					RoomId: "B",
				},
			},
		},
		{
			Id:   "B",
			Desc: "You're standing in a room with a sign that has the letter B on it.",
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
		},
		{
			Id:   "C",
			Desc: "You're standing in a room with a sign that has the letter C on it.",
			Links: []*RoomLink{
				{
					Verb:   "west",
					RoomId: "B",
				},
			},
		},
	}
}

func (w *World) HandleCharacterJoined(character *Character) {
	w.rooms[0].AddCharacter(character)

	character.SendMessage("Welcome to Darkness Falls\n\r")
	character.Room.ShowRoom(character)
	character.Room.ShowOthers(character, fmt.Sprintf("%s appears in a puff of smoke.", character.Name))

	log.Println(fmt.Sprintf("Character login: %s", character.Name))
}

func (w *World) RemoveFromWorld(character *Character) {
	room := character.Room
	room.RemoveCharacter(character)
	room.Show(nil, fmt.Sprintf("%s disappears in a puff of smoke.", character.Name))

	log.Println(fmt.Sprintf("Character logout: %s", character.Name))
}

func (w *World) Broadcast(msg string) {
	for _, player := range w.characters {
		player.SendMessage(msg)
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

func (w *World) HandleCharacterInput(character *Character, input string) {
	room := character.Room
	tokens := strings.Split(input, " ")

	switch tokens[0] {
	case "say":
		msg := strings.Trim(input, "say ")
		character.SendMessage(fmt.Sprintf("You said, '%s'", msg))
		room.ShowOthers(character, fmt.Sprintf("%s said, '%s'", character.Name, msg))

	case "look":
		room.ShowRoom(character)

	case "rename":
		newName := tokens[len(tokens)-1]
		character.Name = newName
		character.SendMessage(fmt.Sprintf("Your name has been changed to %s.", character.Name))

	case "quit":

	default: // direction
		for _, link := range room.Links {
			if link.Verb == input {
				target := w.GetRoomById(link.RoomId)
				if target != nil {
					character.SendMessage(fmt.Sprintf("You travel %s.\r\n", link.Verb))
					room.ShowOthers(character, fmt.Sprintf("%s went %s.", character.Name, link.Verb))
					w.MoveCharacter(character, target)
					character.Room.ShowOthers(character, fmt.Sprintf("%s just came in.", character.Name))
					return
				}
			}
		}
	}
}

func (w *World) MoveCharacter(character *Character, to *Room) {
	character.Room.RemoveCharacter(character)
	to.AddCharacter(character)
	to.ShowRoom(character)
}
