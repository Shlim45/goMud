package main

import (
	"fmt"
	"strings"
)

type RoomLink struct {
	Verb   string
	RoomId string
}

type Direction uint8

const (
	SOUTHWEST = iota
	SOUTH
	SOUTHEAST
	WEST
	EAST
	NORTHWEST
	NORTH
	NORTHEAST
	UP
	DOWN
	OUT
)

func (d Direction) Verb() string {
	switch d {
	case SOUTHWEST:
		return "southwest"
	case SOUTH:
		return "south"
	case SOUTHEAST:
		return "southeast"
	case WEST:
		return "west"
	case EAST:
		return "east"
	case NORTHWEST:
		return "northwest"
	case NORTH:
		return "north"
	case NORTHEAST:
		return "northeast"
	case UP:
		return "up"
	case DOWN:
		return "down"
	case OUT:
		return "out"
	default:
		return "unknown"
	}
}

func VerbToDirection(verb string) Direction {
	switch verb {
	case "southwest":
		return SOUTHWEST
	case "south":
		return SOUTH
	case "southeast":
		return SOUTHEAST
	case "west":
		return WEST
	case "east":
		return EAST
	case "northwest":
		return NORTHWEST
	case "north":
		return NORTH
	case "northeast":
		return NORTHEAST
	case "up":
		return UP
	case "down":
		return DOWN
	case "out":
		return OUT
	default: // TODO(jon): BAD returning OUT for default.
		return OUT
	}
}

type Room struct {
	roomID  string
	Id      string
	Desc    string
	Area    *Area
	Links   []*RoomLink
	Portals []*RoomLink
	Items   []*Item
	Mobs    []*MOB
}

func (r *Room) RoomID() string {
	if r.Area != nil {
		return fmt.Sprintf("%s#%s", r.Area.Name, r.Id)
	}
	return r.Id
}

func (r *Room) AddItem(item *Item) {
	r.Items = append(r.Items, item)
	item.SetOwner(r)
}

func (r *Room) RemoveItem(item *Item) {
	item.SetOwner(nil)

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
		(*item.Owner()).RemoveItem(item)
		item.SetOwner(nil)
	}

	r.Items = append(r.Items, item)
	item.SetOwner(r)
}

func (r *Room) Show(source *MOB, msg string) {
	for _, player := range r.Mobs {
		if player.isPlayer() {
			player.SendMessage(msg, true)
		}
	}
}

func (r *Room) ShowOthers(source *MOB, target *MOB, msg string) {
	for _, player := range r.Mobs {
		if player != nil && player.isPlayer() && player != source && player != target {
			player.SendMessage(msg, true)
		}
	}
}

func (r *Room) ShowRoom(player *MOB) {
	var output strings.Builder

	output.WriteString("[" + CArea(r.Area.Name) + "]\r\n") // area name
	output.WriteString(CNormal("You're " + player.Room.Desc))
	output.WriteString("\r\n")

	numMobs := len(r.Mobs) - 1
	if numMobs > 0 {
		count := 0
		output.WriteString("\r\n")
		for _, other := range r.Mobs {
			if other != player {
				if other.isPlayer() {
					output.WriteString(CFriend(other.Name()))
				} else {
					output.WriteString(CEnemy(other.Name()))
				}
				count++
				if count < numMobs {
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
			output.WriteString(CExit("a " + portal.Verb))
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

	player.SendMessage(output.String(), true)
}

func (r *Room) AddMOB(character *MOB) {
	r.Mobs = append(r.Mobs, character)
	character.Room = r
}

func (r *Room) RemoveMOB(character *MOB) {
	character.Room = nil

	var characters []*MOB
	for _, c := range r.Mobs {
		if c != character {
			characters = append(characters, c)
		}
	}
	r.Mobs = characters
}

func (r *Room) FetchInhabitant(mobName string) *MOB {
	mobName = strings.ToLower(mobName)
	for _, c := range r.Mobs {
		if strings.HasPrefix(strings.ToLower(c.Name()), mobName) {
			return c
		}
	}
	return nil
}

func (r *Room) SaveRoomToDBQuery() string {
	links := []string{"", "", "", "", "", "", "", "", "", "", ""}
	for _, link := range r.Links {
		links[VerbToDirection(link.Verb)] = link.RoomId
	}

	return fmt.Sprintf("INSERT INTO Room VALUES ('%s', '%s', '%s', '%s')",
		r.Id, "Darkness Falls", r.Desc, strings.Join(links, ";"))
}
