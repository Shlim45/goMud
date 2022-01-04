package main

import (
	"fmt"
	"log"
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
			},
		},
	}
}

func (w *World) HandleCharacterJoined(character *Character) {
	w.rooms[0].AddCharacter(character)

	character.SendMessage("Welcome to Darkness Falls\n\r")
	character.SendMessage(character.Room.Desc)
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
	for _, link := range room.Links {
		if link.Verb == input {
			target := w.GetRoomById(link.RoomId)
			if target != nil {
				w.MoveCharacter(character, target)
				return
			}
		}
	}

	character.SendMessage(fmt.Sprintf("You said, '%s'", input))
	for _, other := range room.Characters {
		if other != character {
			other.SendMessage(fmt.Sprintf("%s said, '%s'", character.Name, input))
		}
	}
}

func (w *World) MoveCharacter(character *Character, to *Room) {
	character.Room.RemoveCharacter(character)
	to.AddCharacter(character)
	character.SendMessage(to.Desc)
}
