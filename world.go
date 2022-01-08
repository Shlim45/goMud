package main

import (
	"fmt"
	"log"
	"strings"
)

type World struct {
	characters []*Player
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

func (w *World) HandleCharacterJoined(character *Player) {
	w.rooms[0].AddCharacter(character)

	character.SendMessage("Welcome to Darkness Falls\n\r")
	character.Room.ShowRoom(character)
	character.Room.ShowOthers(character, fmt.Sprintf("%s appears in a puff of smoke.", character.Name))

	log.Println(fmt.Sprintf("Player login: %s", character.Name))
}

func (w *World) RemoveFromWorld(character *Player) {
	room := character.Room
	room.RemoveCharacter(character)
	room.Show(nil, fmt.Sprintf("%s disappears in a puff of smoke.", character.Name))

	log.Println(fmt.Sprintf("Player logout: %s", character.Name))
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

func (w *World) HandlePlayerInput(player *Player, input string) {
	room := player.Room
	tokens := strings.Split(input, " ")

	switch tokens[0] {
	case "say":
		msg := strings.Trim(input, "say ")
		player.SendMessage(fmt.Sprintf("You said, '%s'", msg))
		room.ShowOthers(player, fmt.Sprintf("%s said, '%s'", player.Name, msg))

	case "look":
		room.ShowRoom(player)

	case "rename":
		newName := tokens[len(tokens)-1]
		player.Name = newName
		player.SendMessage(fmt.Sprintf("Your name has been changed to %s.", player.Name))

	case "restat":
		player.Init()
		player.SendMessage("Your stats have been randomized and vitals have been reset to default.")

	case "stats":
		player.SendMessage(fmt.Sprintf("   Att:\t%d Dam:\t%d Eva:\t%d Def:\t%d",
			player.CurPhyStats.Attack, player.CurPhyStats.Damage,
			player.CurPhyStats.Evasion, player.CurPhyStats.Defense))
		player.SendMessage(fmt.Sprintf("MagAtt:\t%d MagDam:\t%d MagEva:\t%d MagDef:\t%d",
			player.CurPhyStats.MagicAttack, player.CurPhyStats.MagicDamage,
			player.CurPhyStats.MagicEvasion, player.CurPhyStats.MagicDefense))

	case "health":
		player.SendMessage(fmt.Sprintf("   Hits: %d/%d     Fat: %d/%d     Pow: %d/%d",
			player.curState().Hits, player.maxState().Hits,
			player.curState().Fat, player.maxState().Fat,
			player.curState().Power, player.maxState().Power))

	case "hit":
		if len(tokens) > 1 {
			target := room.FetchInhabitant(tokens[len(tokens)-1])
			if target != nil {
				player.attackTarget(target)
			} else {
				player.SendMessage("You don't see them here.")
			}
		}

	case "quit":

	default: // direction
		for _, link := range room.Links {
			if link.Verb == input {
				target := w.GetRoomById(link.RoomId)
				if target != nil {
					player.SendMessage(fmt.Sprintf("You travel %s.\r\n", link.Verb))
					room.ShowOthers(player, fmt.Sprintf("%s went %s.", player.Name, link.Verb))
					w.MoveCharacter(player, target)
					player.Room.ShowOthers(player, fmt.Sprintf("%s just came in.", player.Name))
					return
				}
			}
		}
	}
}

func (w *World) MoveCharacter(character *Player, to *Room) {
	character.Room.RemoveCharacter(character)
	to.AddCharacter(character)
	to.ShowRoom(character)
}
