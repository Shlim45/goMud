package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type MudCommand interface {
	ExecuteCmd(m *MOB, input []string, w *World) bool
	Trigger() string
	Timer() time.Duration
	TimerType() TimerType
	CostType() ActionCost
	UseCost() int16
	CheckTimer() bool
}

type Command struct {
	trigger    string
	timer      time.Duration
	timerType  TimerType
	costType   ActionCost
	useCost    int16
	checkTimer bool
}

func (c *Command) Trigger() string {
	return c.trigger
}

func (c *Command) Timer() time.Duration {
	return c.timer
}

func (c *Command) TimerType() TimerType {
	return c.timerType
}

func (c *Command) CostType() ActionCost {
	return c.costType
}

func (c *Command) UseCost() int16 {
	return c.useCost
}

func (c *Command) CheckTimer() bool {
	return c.checkTimer
}

func (c *Command) ExecuteCmd(m *MOB, input []string, w *World) bool {
	success := false
	room := m.Room

	if c.CheckTimer() {
		// if mob still on timer, return false
	}

	// TODO(jon): success = true cases
	switch c.Trigger() {
	case "say":
		msg := strings.Join(input[1:], " ")
		//msg := strings.Trim(input, "say ")
		m.SendMessage(fmt.Sprintf("You said, '%s'", msg), true)
		room.ShowOthers(m, nil, fmt.Sprintf("%s said, '%s'", m.Name(), msg))
		success = true

	case "look":
		room.ShowRoom(m)
		success = true

	case "rename":
		newName := input[len(input)-1]
		m.SetName(newName)
		m.SendMessage(fmt.Sprintf("Your name has been changed to %s.", m.Name()), true)
		if strings.Compare(newName, "Karsus") == 0 {
			m.baseCharStats().setStrength(19)
			m.baseCharStats().setConstitution(15)
			m.baseCharStats().setAgility(20)
			m.baseCharStats().setDexterity(19)
			m.baseCharStats().setIntelligence(17)
			m.baseCharStats().setWisdom(15)
			m.recoverCharStats()
			m.recoverPhyStats()
			m.SendMessage("You shake under the transforming power!", false)
		}
		success = true

	case "reroll":
		m.Init()
		m.SendMessage("Your stats have been randomized and vitals have been reset to default.", true)
		success = true

	case "spawn":
		monster := &MOB{
			name:     "a small dog",
			tickType: TICK_STOP,
		}
		monster.Init()
		monster.basePhyStats().setLevel(5)
		m.Room.AddMOB(monster)
		m.Room.Show(m, fmt.Sprintf("%s appears out of thin air!", CEnemy(monster.Name())))
		success = true

	case "stats":
		var output strings.Builder

		attack := fmt.Sprintf("Attack:    %d/%d", m.curPhyStats().attack(), m.basePhyStats().attack())
		attack = fmt.Sprintf("%-25v", attack)
		output.WriteString(attack)

		mAttack := fmt.Sprintf("Magic Attack:    %d/%d\r\n", m.curPhyStats().magicAttack(), m.basePhyStats().magicAttack())
		output.WriteString(mAttack)

		damage := fmt.Sprintf("Damage:    %d/%d", m.curPhyStats().damage(), m.basePhyStats().damage())
		damage = fmt.Sprintf("%-25v", damage)
		output.WriteString(damage)

		mDamage := fmt.Sprintf("Magic Damage:    %d/%d\r\n", m.curPhyStats().magicDamage(), m.basePhyStats().magicDamage())
		output.WriteString(mDamage)

		evasion := fmt.Sprintf("Evasion:   %d/%d", m.curPhyStats().evasion(), m.basePhyStats().evasion())
		evasion = fmt.Sprintf("%-25v", evasion)
		output.WriteString(evasion)

		mEvasion := fmt.Sprintf("Magic Evasion:   %d/%d\r\n", m.curPhyStats().magicEvasion(), m.basePhyStats().magicEvasion())
		output.WriteString(mEvasion)

		defense := fmt.Sprintf("Defense:   %d/%d", m.curPhyStats().defense(), m.basePhyStats().defense())
		defense = fmt.Sprintf("%-25v", defense)
		output.WriteString(defense)

		mDefense := fmt.Sprintf("Magic Defense:   %d/%d\r\n", m.curPhyStats().magicDefense(), m.basePhyStats().magicDefense())
		output.WriteString(mDefense)

		m.SendMessage(output.String(), true)
		success = true

	case "health":
		m.SendMessage(fmt.Sprintf("   Hits: %d/%d     Fat: %d/%d     Pow: %d/%d",
			m.curState().Hits, m.maxState().Hits,
			m.curState().Fat, m.maxState().Fat,
			m.curState().Power, m.maxState().Power), true)
		success = true

	case "release":
		if m.curState().Alive {
			m.SendMessage("You must be dead to release your corpse.", true)
			return success
		}
		m.releaseCorpse(w)
		success = true

	case "hit":
		if len(input) > 1 {
			target := room.FetchInhabitant(input[len(input)-1])
			if target != nil {
				m.attackTarget(target)
				success = true
			} else {
				m.SendMessage("You don't see them here.", true)
			}
		} else {
			target := m.Victim
			if target != nil {
				m.attackTarget(target)
				success = true
			} else {
				m.SendMessage("Hit who?  You must specify a target.", true)
			}
		}

	case "information":
		var output strings.Builder

		name := fmt.Sprintf("Name:     %s", m.Name())
		name = fmt.Sprintf("%-25v", name)
		output.WriteString(name)

		pClass := fmt.Sprintf("Class:    %s\r\n", "Generic Class")
		output.WriteString(pClass)

		level := fmt.Sprintf("Level:    %d", m.curPhyStats().level())
		level = fmt.Sprintf("%-25v", level)
		output.WriteString(level)

		exp := fmt.Sprintf("Experience:    %d\r\n", m.Experience)
		output.WriteString(exp)

		output.WriteString(Yellow("\r\n                  Score   Bonus\r\n"))
		output.WriteString("                  -----   -----\r\n")
		output.WriteString(fmt.Sprintf("Strength:         %2v\r\n", m.curCharStats().strength()))
		output.WriteString(fmt.Sprintf("Constitution:     %2v\r\n", m.curCharStats().constitution()))
		output.WriteString(fmt.Sprintf("Agility:          %2v\r\n", m.curCharStats().agility()))
		output.WriteString(fmt.Sprintf("Dexterity:        %2v\r\n", m.curCharStats().dexterity()))
		output.WriteString(fmt.Sprintf("Intelligence:     %2v\r\n", m.curCharStats().intelligence()))
		output.WriteString(fmt.Sprintf("Wisdom:           %2v\r\n", m.curCharStats().wisdom()))

		output.WriteString("\r\nTry STATS or HEALTH commands.")

		m.SendMessage(output.String(), true)
		success = true

	case "quit":
		m.SendMessage("Not working yet, Ctrl+] to quit from telnet prompt.", true)

	case "inventory":
		m.ShowInventory()
		success = true

	case "wealth":
		m.SendMessage(fmt.Sprintf("You are carrying %d silver coins.", m.Coins), true)
		success = true

	case "create":
		if len(input) < 2 {
			m.SendMessage("Create what?  Syntax: CREATE <ARTICLE OR .> <ITEM NAME> <KEYWORD>", true)
			return success
		}

		var newItem Item

		if n, err := strconv.Atoi(input[1]); err == nil {
			if n <= 0 {
				m.SendMessage("You must specify a number greater than 0 to create coins!", true)
				return success
			}

			newItem = Item{
				article:  "a",
				name:     "pile of silver coins",
				keyword:  "coins",
				owner:    nil,
				value:    uint64(n),
				itemType: TYPE_COINS,
			}

		} else {
			itemArticle := input[1]
			itemName := strings.Join(input[2:len(input)-1], " ")
			itemKeyword := input[len(input)-1]

			if strings.Compare(itemArticle, ".") == 0 {
				itemArticle = ""
			}

			newItem = Item{
				keyword:  itemKeyword,
				article:  itemArticle,
				name:     itemName,
				owner:    nil,
				value:    0,
				itemType: TYPE_GENERIC,
			}
		}

		m.Room.AddItem(&newItem)
		m.Room.Show(nil, fmt.Sprintf("\r\n%s falls from the sky!", newItem.Name()))
		success = true

	case "get":
		if len(input) < 2 {
			m.SendMessage("Get what?  Syntax: GET <ITEM>", true)
			return success
		}

		targetItem := input[1]
		for _, item := range m.Room.Items {
			if strings.HasPrefix(item.Keyword(), targetItem) {
				if item.ItemType() == ItemType(TYPE_COINS) {
					numCoins := item.Value()
					m.Room.RemoveItem(item)
					m.Coins += numCoins
					m.SendMessage(fmt.Sprintf("You pick up %d silver coins.", numCoins), true)
				} else {
					m.MoveItemTo(item)
					m.SendMessage(fmt.Sprintf("You pick up %s.", item.FullName()), true)
				}
				m.Room.ShowOthers(m, nil, fmt.Sprintf("%s picks up %s.", m.Name(), item.FullName()))
				success = true
				break
			}
		}

		if !success {
			m.SendMessage(fmt.Sprintf("You don't see a '%s' here.", targetItem), true)
			return success
		}

	case "drop":
		if len(input) < 2 {
			m.SendMessage("Drop what?  Syntax: DROP <ITEM>", true)
			return success
		}

		if n, err := strconv.Atoi(input[1]); err == nil {
			if n <= 0 {
				m.SendMessage("You must specify a number greater than 0 to drop coins!", true)
				return success
			}

			if m.Coins < uint64(n) {
				m.SendMessage("You aren't carrying enough silver!", true)
				return success
			}

			coinsItem := Item{
				article: "a",
				name:    "pile of silver coins",
				keyword: "coins",
				owner:   nil,
				value:   uint64(n),
			}

			m.Coins -= uint64(n)
			m.Room.MoveItemTo(&coinsItem)
			m.SendMessage(fmt.Sprintf("You drop %d silver coins.", n), true)
			m.Room.ShowOthers(m, nil, fmt.Sprintf("%s drops %s.", m.Name(), coinsItem.FullName()))
			success = true
			break
		}

		targetItem := input[1]
		if !success {
			for _, item := range m.Inventory() {
				if strings.HasPrefix(item.Keyword(), targetItem) {
					m.Room.MoveItemTo(item)
					m.SendMessage(fmt.Sprintf("You drop %s.", item.FullName()), true)
					m.Room.ShowOthers(m, nil, fmt.Sprintf("%s drops %s.", m.Name(), item.FullName()))
					success = true
					break
				}
			}
		}

		if !success {
			m.SendMessage(fmt.Sprintf("You aren't carrying a '%s'.", targetItem), true)
			return success
		}

	case "give":
		if len(input) < 3 {
			m.SendMessage("Give what to whom?  Syntax: GIVE <ITEM> <TARGET>", true)
			return success
		}

		itemName := input[1]
		targetName := input[2]

		targetMob := m.Room.FetchInhabitant(targetName)
		if targetMob == nil {
			m.SendMessage("You don't see them here.", true)
			return success
		}

		if n, err := strconv.Atoi(input[1]); err == nil {
			numCoins := uint64(n)

			if numCoins <= 0 {
				m.SendMessage("You must specify a number greater than 0 to give coins!", true)
				return success
			}
			if m.Coins < numCoins {
				m.SendMessage("You aren't carrying enough silver!", true)
				return success
			}

			m.Coins -= numCoins
			targetMob.Coins += numCoins

			m.SendMessage(fmt.Sprintf("You give %d silver coins to %s.",
				numCoins, targetMob.Name()), true)
			targetMob.SendMessage(fmt.Sprintf("%s gives you %d silver coins.",
				m.Name(), numCoins), true)
			m.Room.ShowOthers(m, targetMob, fmt.Sprintf("%s gives %s some silver coins.",
				m.Name(), targetMob.Name()))
			success = true
		}

		if !success {
			for _, item := range m.Inventory() {
				if strings.HasPrefix(item.Keyword(), itemName) {
					targetMob.MoveItemTo(item)
					m.SendMessage(fmt.Sprintf("You give %s to %s.",
						item.FullName(), targetMob.Name()), true)
					targetMob.SendMessage(fmt.Sprintf("%s gives you %s.",
						m.Name(), item.FullName()), true)
					m.Room.ShowOthers(m, targetMob, fmt.Sprintf("%s gives %s %s.",
						m.Name(), targetMob.Name(), item.FullName()))
					success = true
					break
				}
			}
		}

		if !success {
			m.SendMessage(fmt.Sprintf("You aren't carrying a '%s'.", itemName), true)
			return success
		}

	case "go":
		if len(input) < 2 {
			m.SendMessage("Go where?  Syntax: GO <EXIT>", true)
			return success
		}

		if m.curState().fat() < 2 {
			m.SendMessage("You are too tired!", true)
			return success
		}

		//targetPortal := strings.Join(input[1:], " ")
		targetPortal := input[1]

		for _, portal := range room.Portals {
			if portal.Verb == targetPortal || strings.HasPrefix(portal.Verb, targetPortal) {
				target := w.GetRoomById(portal.RoomId)
				if target != nil {
					m.SendMessage(fmt.Sprintf("You travel into a %s.", portal.Verb), true)
					room.ShowOthers(m, nil, fmt.Sprintf("%s went into a %s.", m.Name(), portal.Verb))
					w.MoveCharacter(m, target)
					m.Room.ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
					success = true
					break
				}
			}
		}

		if !success {
			for _, link := range room.Links {
				if link.Verb == targetPortal || strings.HasPrefix(link.Verb, targetPortal) {
					target := w.GetRoomById(link.RoomId)
					if target != nil {
						m.SendMessage(fmt.Sprintf("You travel %s.", link.Verb), true)
						room.ShowOthers(m, nil, fmt.Sprintf("%s went %s.", m.Name(), link.Verb))
						w.MoveCharacter(m, target)
						m.Room.ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
						success = true
						break
					}
				}
			}
		}

		if !success {
			m.SendMessage(fmt.Sprintf("There isn't a '%s' here.", targetPortal), true)
			return success
		}

	default:
		log.Panicf("Command with trigger '%s' not found.", c.Trigger())
		return success
	}

	if success {
		// handle timer and usage cost
		if c.CostType() == COST_FAT {
			m.adjFat(-c.UseCost(), m.maxState().fat())
		}
	}
	return success
}
