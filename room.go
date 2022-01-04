package main

import (
	"strings"
)

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

func (r *Room) Show(source *Character, msg string) {
	for _, player := range r.Characters {
		player.SendMessage(msg)
	}
}

func (r *Room) ShowOthers(source *Character, msg string) {
	for _, player := range r.Characters {
		if player != nil && player != source {
			player.SendMessage(msg)
		}
	}
}

func (r *Room) ShowRoom(character *Character) {
	//character.SendMessage(character.Room.Desc)
	var output strings.Builder

	output.WriteString("[Darkness Falls]\r\n") // area name
	output.WriteString(character.Room.Desc)
	output.WriteString("\r\n")

	numOthers := len(r.Characters) - 1
	if numOthers > 0 {
		count := 0
		output.WriteString("\r\nAlso there is ")
		for _, other := range r.Characters {
			if other != character {
				output.WriteString(other.Name)
				count++
				if count < numOthers {
					output.WriteString(", ")
				} else {
					output.WriteString(".\r\n")
				}
			}
		}
	}

	numExits := len(r.Links)
	if numExits > 0 {
		count := 0
		output.WriteString("\r\nObvious Exits: ")
		for _, link := range r.Links {
			output.WriteString(link.Verb)
			count++
			if count < numExits {
				output.WriteString(", ")
			} else {
				output.WriteString(".\r\n")
			}
		}
	}

	character.SendMessage(output.String())
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
