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
			Portals: []*RoomLink{},
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
			Portals: []*RoomLink{
				{
					Verb:   "gate",
					RoomId: "D",
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
			Portals: []*RoomLink{},
		},
		{
			Id:    "D",
			Desc:  "You're standing in a room hidden behind a gate.  There is a sign that has the letter D on it.",
			Links: []*RoomLink{},
			Portals: []*RoomLink{
				{
					Verb:   "gate",
					RoomId: "B",
				},
			},
		},
	}
}

func (w *World) HandleCharacterJoined(character *Player) {
	w.rooms[0].AddCharacter(character)

	character.SendMessage("Welcome to Darkness Falls\n\r", true)
	character.Room.ShowRoom(character)
	character.Room.ShowOthers(character, nil, fmt.Sprintf("%s appears in a puff of smoke.", character.Name))

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

func (w *World) MoveCharacter(character *Player, to *Room) {
	character.Room.RemoveCharacter(character)
	to.AddCharacter(character)
	to.ShowRoom(character)
}

func (w *World) HandlePlayerInput(player *Player, input string) {
	room := player.Room
	tokens := strings.Split(input, " ")

	switch tokens[0] {
	case "say":
		msg := strings.Trim(input, "say ")
		player.SendMessage(fmt.Sprintf("You said, '%s'", msg), true)
		room.ShowOthers(player, nil, fmt.Sprintf("%s said, '%s'", player.Name, msg))

	case "look":
		room.ShowRoom(player)

	case "rename":
		newName := tokens[len(tokens)-1]
		player.Name = newName
		player.SendMessage(fmt.Sprintf("Your name has been changed to %s.", player.Name), true)

	case "restat":
		player.Init()
		player.SendMessage("Your stats have been randomized and vitals have been reset to default.", true)

	case "stats":
		var output strings.Builder

		attack := fmt.Sprintf("Attack:    %d/%d", player.curPhyStats().attack(), player.basePhyStats().attack())
		attack = fmt.Sprintf("%-25v", attack)
		output.WriteString(attack)

		mAttack := fmt.Sprintf("Magic Attack:    %d/%d\r\n", player.curPhyStats().magicAttack(), player.basePhyStats().magicAttack())
		output.WriteString(mAttack)

		damage := fmt.Sprintf("Damage:    %d/%d", player.curPhyStats().damage(), player.basePhyStats().damage())
		damage = fmt.Sprintf("%-25v", damage)
		output.WriteString(damage)

		mDamage := fmt.Sprintf("Magic Damage:    %d/%d\r\n", player.curPhyStats().magicDamage(), player.basePhyStats().magicDamage())
		output.WriteString(mDamage)

		evasion := fmt.Sprintf("Evasion:   %d/%d", player.curPhyStats().evasion(), player.basePhyStats().evasion())
		evasion = fmt.Sprintf("%-25v", evasion)
		output.WriteString(evasion)

		mEvasion := fmt.Sprintf("Magic Evasion:   %d/%d\r\n", player.curPhyStats().magicEvasion(), player.basePhyStats().magicEvasion())
		output.WriteString(mEvasion)

		defense := fmt.Sprintf("Defense:   %d/%d", player.curPhyStats().defense(), player.basePhyStats().defense())
		defense = fmt.Sprintf("%-25v", defense)
		output.WriteString(defense)

		mDefense := fmt.Sprintf("Magic Defense:   %d/%d\r\n", player.curPhyStats().magicDefense(), player.basePhyStats().magicDefense())
		output.WriteString(mDefense)

		player.SendMessage(output.String(), true)

	case "health":
		player.SendMessage(fmt.Sprintf("   Hits: %d/%d     Fat: %d/%d     Pow: %d/%d",
			player.curState().Hits, player.maxState().Hits,
			player.curState().Fat, player.maxState().Fat,
			player.curState().Power, player.maxState().Power), true)

	case "recall":
		if player.curState().Alive {
			player.SendMessage("You must be dead to recall your corpse.", true)
			return
		}
		player.recallCorpse(w)

	case "hit":
		if len(tokens) > 1 {
			target := room.FetchInhabitant(tokens[len(tokens)-1])
			if target != nil {
				player.attackTarget(target)
			} else {
				player.SendMessage("You don't see them here.", true)
			}
		}

	case "info":
		var output strings.Builder

		name := fmt.Sprintf("Name:     %s", player.Name)
		name = fmt.Sprintf("%-25v", name)
		output.WriteString(name)

		pClass := fmt.Sprintf("Class:    %s\r\n", "Generic Class")
		output.WriteString(pClass)

		level := fmt.Sprintf("Level:    %d", player.curPhyStats().level())
		level = fmt.Sprintf("%-25v", level)
		output.WriteString(level)

		exp := fmt.Sprintf("Experience:    %d\r\n", player.Experience)
		output.WriteString(exp)

		output.WriteString(Yellow("\r\n                  Score   Bonus\r\n"))
		output.WriteString("                  -----   -----\r\n")
		output.WriteString(fmt.Sprintf("Strength:         %2v\r\n", player.curCharStats().strength()))
		output.WriteString(fmt.Sprintf("Constitution:     %2v\r\n", player.curCharStats().constitution()))
		output.WriteString(fmt.Sprintf("Agility:          %2v\r\n", player.curCharStats().agility()))
		output.WriteString(fmt.Sprintf("Dexterity:        %2v\r\n", player.curCharStats().dexterity()))
		output.WriteString(fmt.Sprintf("Intelligence:     %2v\r\n", player.curCharStats().intelligence()))
		output.WriteString(fmt.Sprintf("Wisdom:           %2v\r\n", player.curCharStats().wisdom()))

		output.WriteString("\r\nTry STATS or HEALTH commands.")

		player.SendMessage(output.String(), true)

	case "quit":
		player.SendMessage("Not working yet, Ctrl+] to quit from telnet prompt.", true)

	case "go":
		if len(tokens) < 2 {
			player.SendMessage("Go where?  Syntax: GO <EXIT>", true)
			return
		}

		targetPortal := strings.Join(tokens[1:], " ")

		for _, portal := range room.Portals {
			if portal.Verb == targetPortal || strings.HasPrefix(portal.Verb, targetPortal) {
				target := w.GetRoomById(portal.RoomId)
				if target != nil {
					player.SendMessage(fmt.Sprintf("You travel into a %s.", portal.Verb), true)
					room.ShowOthers(player, nil, fmt.Sprintf("%s went into a %s.", player.Name, portal.Verb))
					w.MoveCharacter(player, target)
					player.Room.ShowOthers(player, nil, fmt.Sprintf("%s just came in.", player.Name))
					return
				}
			}
		}

		for _, link := range room.Links {
			if link.Verb == targetPortal || strings.HasPrefix(link.Verb, targetPortal) {
				target := w.GetRoomById(link.RoomId)
				if target != nil {
					player.SendMessage(fmt.Sprintf("You travel %s.", link.Verb), true)
					room.ShowOthers(player, nil, fmt.Sprintf("%s went %s.", player.Name, link.Verb))
					w.MoveCharacter(player, target)
					player.Room.ShowOthers(player, nil, fmt.Sprintf("%s just came in.", player.Name))
					return
				}
			}
		}

		player.SendMessage(fmt.Sprintf("There isn't a '%s' here.", targetPortal), true)

	default: // direction
		if !player.curState().Alive {
			player.SendMessage("You must be alive to do that!", true)
			return
		}

		for _, link := range room.Links {
			if link.Verb == input || strings.HasPrefix(link.Verb, input) {
				target := w.GetRoomById(link.RoomId)
				if target != nil {
					player.SendMessage(fmt.Sprintf("You travel %s.", link.Verb), true)
					room.ShowOthers(player, nil, fmt.Sprintf("%s went %s.", player.Name, link.Verb))
					w.MoveCharacter(player, target)
					player.Room.ShowOthers(player, nil, fmt.Sprintf("%s just came in.", player.Name))
					return
				}
			}
		}

		player.SendMessage(fmt.Sprintf("Huh?  Command or Exit '%s' not found.", input), true)
	}
}
