package main

import (
	"strings"
)

type RoomLink struct {
	Verb   string
	RoomId string
}

type Room struct {
	Id      string
	Desc    string
	Links   []*RoomLink
	Portals []*RoomLink

	Characters []*Player
}

func (r *Room) Show(source *Player, msg string) {
	for _, player := range r.Characters {
		player.SendMessage(msg, true)
	}
}

func (r *Room) ShowOthers(source *Player, msg string) {
	for _, player := range r.Characters {
		if player != nil && player != source {
			player.SendMessage(msg, true)
		}
	}
}

func (r *Room) ShowRoom(character *Player) {
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

	numPortals := len(r.Portals)
	if numPortals > 0 {
		count := 0
		output.WriteString("You also see ")
		for _, portal := range r.Portals {
			output.WriteString("a " + portal.Verb)
			count++
			if count < numPortals {
				output.WriteString(", ")
			} else {
				output.WriteString(".\r\n")
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

	character.SendMessage(output.String(), true)
}

func (r *Room) AddCharacter(character *Player) {
	r.Characters = append(r.Characters, character)
	character.Room = r
}

func (r *Room) RemoveCharacter(character *Player) {
	character.Room = nil

	var characters []*Player
	for _, c := range r.Characters {
		if c != character {
			characters = append(characters, c)
		}
	}
	r.Characters = characters
}

func (r *Room) FetchInhabitant(mobName string) *Player {
	mobName = strings.ToLower(mobName)
	for _, c := range r.Characters {
		if strings.HasPrefix(strings.ToLower(c.Name), mobName) {
			return c
		}
	}
	return nil
}
