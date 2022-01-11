package main

import (
	"strings"
)

type RoomLink struct {
	Verb   string
	RoomId string
}

type Room struct {
	Id         string
	Desc       string
	Links      []*RoomLink
	Portals    []*RoomLink
	Items      []*Item
	Characters []*Player
}

func (r *Room) AddItem(item *Item) {
	r.Items = append(r.Items, item)
	item.SetOwner(nil)
	item.SetLocation(r)
}

func (r *Room) RemoveItem(item *Item) {
	item.SetOwner(nil)
	item.SetLocation(nil)

	var items []*Item
	for _, i := range r.Items {
		if i != item {
			items = append(items, i)
		}
	}
	r.Items = items
}

func (r *Room) MoveItemTo(item *Item) {
	if item.Owner() != nil {
		item.Owner().RemoveItem(item)
		item.SetOwner(nil)
	}

	if item.Location() != nil {
		item.Location().RemoveItem(item)
		item.SetLocation(nil)
	}

	r.Items = append(r.Items, item)
	item.SetLocation(r)
}

func (r *Room) Show(source *Player, msg string) {
	for _, player := range r.Characters {
		player.SendMessage(msg, true)
	}
}

func (r *Room) ShowOthers(source *Player, target *Player, msg string) {
	for _, player := range r.Characters {
		if player != nil && player != source && player != target {
			player.SendMessage(msg, true)
		}
	}
}

func (r *Room) ShowRoom(character *Player) {
	//character.SendMessage(character.Room.Desc)
	var output strings.Builder

	output.WriteString("[" + CArea("Darkness Falls") + "]\r\n") // area name
	output.WriteString(CNormal(character.Room.Desc))
	output.WriteString("\r\n")

	numOthers := len(r.Characters) - 1
	if numOthers > 0 {
		count := 0
		output.WriteString("\r\n")
		for _, other := range r.Characters {
			if other != character {
				output.WriteString(CFriend(other.Name()))
				count++
				if count < numOthers {
					output.WriteString(", ")
				} else {
					output.WriteString(" is also here.\r\n")
				}
			}
		}
	}

	numPortals := len(r.Portals)
	if numPortals > 0 {
		count := 0
		output.WriteString("You also see ")
		for _, portal := range r.Portals {
			output.WriteString("a " + CExit(portal.Verb))
			count++
			if count < numPortals {
				output.WriteString(", ")
			} else {
				output.WriteString(".\r\n")
			}
		}
	}

	numItems := len(r.Items)
	if numItems > 0 {
		count := 0
		output.WriteString("You also see ")
		for _, item := range r.Items {
			output.WriteString(CItem(item.FullName()))
			count++
			if count < numItems {
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
		if strings.HasPrefix(strings.ToLower(c.Name()), mobName) {
			return c
		}
	}
	return nil
}
