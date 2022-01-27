package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type MudCommand interface {
	ExecuteCmd(m *Player, input []string, w *World) bool
	Trigger() string
	Timer() time.Duration
	TimerType() TimerType
	CostType() ActionCost
	UseCost() int16
	CheckTimer() bool
	Clearance() Security
}

type Command struct {
	trigger    string
	timer      time.Duration
	timerType  TimerType
	costType   ActionCost
	useCost    int16
	checkTimer bool
	security   Security
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

func (c *Command) ExecuteCmd(m *Player, input []string, library *MudLib) bool {
	success := false
	room := m.Room()

	if c.CheckTimer() {
		// if mob still on timer, return false
	}

	switch c.Trigger() {
	case "say":
		msg := strings.Join(input[1:], " ")

		m.SendMessage(fmt.Sprintf("You said, '%s'", msg), true)
		room.ShowOthers(m, nil, fmt.Sprintf("%s said, '%s'", m.Name(), msg))
		success = true

	case "look":
		room.ShowRoom(m)
		success = true

	case "rename":
		newName := input[len(input)-1]
		_, exists := library.world.characters[newName]
		if !exists {
			delete(library.world.characters, m.Name())
			m.SetName(newName)
			library.world.characters[m.Name()] = m
			m.SendMessage(fmt.Sprintf("Your name has been changed to %s.", m.Name()), true)
			success = true
		} else {
			m.SendMessage("That name is already in use.", true)
		}

	case "reroll":
		m.Init(library)
		m.SendMessage("Your stats have been randomized and vitals have been reset to default.", true)
		success = true

	case "stats":
		var output strings.Builder

		attack := fmt.Sprintf("Attack:    %d/%d", m.CurPhyStats().attack(), m.BasePhyStats().attack())
		attack = fmt.Sprintf("%-25v", attack)
		output.WriteString(attack)

		mAttack := fmt.Sprintf("Magic Attack:    %d/%d\r\n", m.CurPhyStats().magicAttack(), m.BasePhyStats().magicAttack())
		output.WriteString(mAttack)

		damage := fmt.Sprintf("Damage:    %d/%d", m.CurPhyStats().damage(), m.BasePhyStats().damage())
		damage = fmt.Sprintf("%-25v", damage)
		output.WriteString(damage)

		mDamage := fmt.Sprintf("Magic Damage:    %d/%d\r\n", m.CurPhyStats().magicDamage(), m.BasePhyStats().magicDamage())
		output.WriteString(mDamage)

		evasion := fmt.Sprintf("Evasion:   %d/%d", m.CurPhyStats().evasion(), m.BasePhyStats().evasion())
		evasion = fmt.Sprintf("%-25v", evasion)
		output.WriteString(evasion)

		mEvasion := fmt.Sprintf("Magic Evasion:   %d/%d\r\n", m.CurPhyStats().magicEvasion(), m.BasePhyStats().magicEvasion())
		output.WriteString(mEvasion)

		defense := fmt.Sprintf("Defense:   %d/%d", m.CurPhyStats().defense(), m.BasePhyStats().defense())
		defense = fmt.Sprintf("%-25v", defense)
		output.WriteString(defense)

		mDefense := fmt.Sprintf("Magic Defense:   %d/%d\r\n", m.CurPhyStats().magicDefense(), m.BasePhyStats().magicDefense())
		output.WriteString(mDefense)

		m.SendMessage(output.String(), true)
		success = true

	case "health":
		m.SendMessage(fmt.Sprintf("   Hits: %d/%d     Fat: %d/%d     Pow: %d/%d",
			m.CurState().Hits, m.MaxState().Hits,
			m.CurState().Fat, m.MaxState().Fat,
			m.CurState().Power, m.MaxState().Power), true)
		success = true

	case "release":
		if m.CurState().Alive {
			m.SendMessage("You must be dead to release your corpse.", true)
			return success
		}
		m.releaseCorpse(library.world)
		success = true

	case "hit":
		if len(input) > 1 {
			target := room.FetchInhabitant(input[len(input)-1])
			if target != nil {
				m.AttackTarget(target)
				success = true
			} else {
				m.SendMessage("You don't see them here.", true)
			}
		} else {
			target := m.Victim()
			if target != nil {
				if target.Room() != m.Room() {
					m.SendMessage("You don't see them here.", true)
					return success
				}
				m.AttackTarget(target)
				success = true
			} else {
				m.SendMessage("Hit who?  You must specify a target.", true)
			}
		}

	case "information":
		var output strings.Builder

		name := fmt.Sprintf("Name:     %s", m.Name())
		name = fmt.Sprintf("%-40v", name)
		output.WriteString(name)

		pClass := fmt.Sprintf("Class:    %s\r\n", m.CurCharStats().CurrentClass().Name())
		output.WriteString(pClass)

		level := fmt.Sprintf("Level:    %d", m.CurPhyStats().level())
		level = fmt.Sprintf("%-40v", level)
		output.WriteString(level)

		exp := fmt.Sprintf("Experience:    %d\r\n", m.Experience)
		output.WriteString(exp)

		var realmTitle string
		if m.RealmPoints > 0 {
			realmTitle = "Dark Acolyte"
		} else {
			realmTitle = ""
		}
		realmRank := fmt.Sprintf("Realm Title:    %s", realmTitle)
		realmRank = fmt.Sprintf("%-40v", realmRank)
		output.WriteString(realmRank)

		rp := fmt.Sprintf("Realm Points:  %d\r\n", m.RealmPoints)
		output.WriteString(rp)

		output.WriteString(Yellow("\r\n                  Score   Bonus\r\n"))
		output.WriteString("                  -----   -----\r\n")

		for stat, value := range m.CurCharStats().Stats() {
			bonus := m.CurCharStats().StatBonus(uint8(stat))
			output.WriteString(fmt.Sprintf("%-18v%2v (%2v%%)\r\n",
				StatToString(uint8(stat))+":", value, bonus))
		}

		output.WriteString("\r\nTry STATS or HEALTH commands.")

		m.SendMessage(output.String(), true)
		success = true

	case "quit":
		m.SendMessage("Are you sure you want to quit? (y/N)", true)
		m.User.Session.SetStatus(QUIT)

	case "inventory":
		m.ShowInventory()
		success = true

	case "wealth":
		m.SendMessage(fmt.Sprintf("You are carrying %d silver coins.", m.Coins), true)
		success = true

	case "get":
		if len(input) < 2 {
			m.SendMessage("Get what?  Syntax: GET <ITEM>", true)
			return success
		}

		targetItem := input[1]
		item := m.Room().FindItem(targetItem)
		if item != nil {
			if item.ItemType() == ItemType(TYPE_COINS) {
				numCoins := item.Value()
				m.Room().RemoveItem(item)
				m.AdjCoins(int64(numCoins))
				m.SendMessage(fmt.Sprintf("You pick up %d silver coins.", numCoins), true)
			} else {
				m.MoveItemTo(item)
				m.SendMessage(fmt.Sprintf("You pick up %s.", item.FullName()), true)
			}
			m.Room().ShowOthers(m, nil, fmt.Sprintf("%s picks up %s.", m.Name(), item.FullName()))
			success = true
		} else {
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

			if m.Coins() < uint64(n) {
				m.SendMessage("You aren't carrying enough silver!", true)
				return success
			}

			coinsItem := Item{
				article: "a",
				name:    "pile of silver coins",
				keyword: "pile",
				owner:   nil,
				value:   uint64(n),
			}

			m.AdjCoins(-int64(n))
			m.Room().MoveItemTo(&coinsItem)
			m.SendMessage(fmt.Sprintf("You drop %d silver coins.", n), true)
			m.Room().ShowOthers(m, nil, fmt.Sprintf("%s drops %s.", m.Name(), coinsItem.FullName()))
			success = true
			break
		}

		targetItem := input[1]
		if !success {
			for _, item := range m.Inventory() {
				if strings.HasPrefix(item.Keyword(), targetItem) {
					m.Room().MoveItemTo(item)
					m.SendMessage(fmt.Sprintf("You drop %s.", item.FullName()), true)
					m.Room().ShowOthers(m, nil, fmt.Sprintf("%s drops %s.", m.Name(), item.FullName()))
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

		targetMob := m.Room().FetchInhabitant(targetName)
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
			if m.Coins() < numCoins {
				m.SendMessage("You aren't carrying enough silver!", true)
				return success
			}

			m.AdjCoins(-int64(numCoins))
			targetMob.AdjCoins(int64(numCoins))

			m.SendMessage(fmt.Sprintf("You give %d silver coins to %s.",
				numCoins, targetMob.Name()), true)
			targetMob.SendMessage(fmt.Sprintf("%s gives you %d silver coins.",
				m.Name(), numCoins), true)
			m.Room().ShowOthers(m, targetMob, fmt.Sprintf("%s gives %s some silver coins.",
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
					m.Room().ShowOthers(m, targetMob, fmt.Sprintf("%s gives %s %s.",
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

		if m.CurState().fat() < 2 {
			m.SendMessage("You are too tired!", true)
			return success
		}

		targetPortal := input[1]

		port := m.Room().FindExit(targetPortal)
		if port != nil {
			destRoom := port.DestRoom()
			if destRoom != nil {
				m.WalkThrough(port)
				success = true
				break
			} else {
				m.SendMessage("That exit is broken... nowhere to go!", true)
			}
		} else {
			for _, link := range room.Links {
				if link.Verb == targetPortal || strings.HasPrefix(link.Verb, targetPortal) {
					target := library.world.GetRoomById(link.RoomId)
					if target != nil {
						m.Walk(target, link.Verb)
						success = true
						break
					}
				}
			}

			if !success {
				m.SendMessage(fmt.Sprintf("There isn't a '%s' here.", targetPortal), true)
				return success
			}
		}

	case "*goto":
		if m.SecClearance >= c.security {
			if len(input) < 2 {
				m.SendMessage("GoTo where?  Syntax: *GOTO <ROOMID> or <PLAYERNAME>", true)
				return success
			}

			targetDest := strings.Join(input[1:], " ")
			destRoom := library.world.GetRoomById(targetDest)

			if destRoom == nil {
				destPlayer := library.world.characters[targetDest]
				if destPlayer != nil {
					destRoom = destPlayer.Room()
				}
			}

			if destRoom != nil {
				library.world.MoveMob(m, destRoom)
				success = true
			} else {
				m.SendMessage(fmt.Sprintf("Unable to locate destination '%s' as RoomID or Player.", targetDest), true)
			}
		}

	case "*spawn":
		if m.SecClearance >= c.security {
			monster := &Monster{
				name:     "a small dog",
				tickType: TickStop,
			}
			monster.Init(library)
			monster.BasePhyStats().setLevel(5)
			m.Room().AddMOB(monster)
			m.Room().Show(m, fmt.Sprintf("%s appears out of thin air!", CEnemy(monster.Name())))
			success = true
		}

	case "*create":
		if m.SecClearance >= c.security {

			if len(input) < 2 {
				m.SendMessage("Create what?  Syntax: CREATE <OBJECT TYPE> <ARTICLE OR .> <OBJECT NAME> <KEYWORD> [(<DESTINATION ROOMID>)]", true)
				return success
			}

			var newItem Item

			if n, err := strconv.Atoi(input[1]); err == nil {
				// creating silver coins
				if n <= 0 {
					m.SendMessage("You must specify a number greater than 0 to create coins!", true)
					return success
				}

				newItem = Item{
					article:  "a",
					name:     "pile of silver coins",
					keyword:  "pile",
					owner:    nil,
					value:    uint64(n),
					itemType: TYPE_COINS,
				}

				m.Room().AddItem(&newItem)
				m.Room().Show(nil, fmt.Sprintf("\r\n%s falls from the sky!", newItem.Name()))
				success = true

			} else {
				// creating an object
				itemType := strings.ToLower(input[1])
				createItem := strings.Compare(itemType, "item") == 0
				createExit := strings.Compare(itemType, "exit") == 0

				if !createItem && !createExit {
					m.SendMessage("Create what?  Can create ITEM or EXIT.", true)
					return success
				}

				if createItem {
					if len(input) < 5 {
						m.SendMessage("You failed to specify all required fields."+
							"  Syntax: CREATE ITEM <ARTICLE OR .> <NAME> <KEYWORD>", true)
						return success
					}
					itemArticle := input[2]
					itemName := strings.Join(input[3:len(input)-1], " ")
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

					m.Room().AddItem(&newItem)
					m.Room().Show(nil, fmt.Sprintf("\r\n%s falls from the sky!", newItem.Name()))
					success = true
				} else if createExit {
					if len(input) < 6 {
						m.SendMessage("You failed to specify all required fields."+
							"  Syntax: CREATE EXIT <ARTICLE OR .> <NAME> <KEYWORD> (<DESTINATION ROOMID>)", true)
						return success
					}

					joinedInput := strings.Join(input, " ")

					namingInput := before(joinedInput, " (")
					namingTokens := strings.Split(namingInput, " ")

					destRoomID := after(joinedInput, "(")
					destRoomID = strings.TrimSuffix(destRoomID, ")")
					destRoom := library.world.GetRoomById(destRoomID)
					if destRoom == nil {
						m.SendMessage(fmt.Sprintf("Invalid RoomID '%s'.", destRoomID), true)
						return success
					}

					exitArticle := namingTokens[2]
					exitName := strings.Join(namingTokens[3:len(namingTokens)-1], " ")
					exitKeyword := namingTokens[len(namingTokens)-1]

					if strings.Compare(exitArticle, ".") == 0 {
						exitArticle = ""
					}
					newExit := Portal{
						room:     m.Room(),
						destRoom: destRoom,
						article:  exitArticle,
						name:     exitName,
						keyword:  exitKeyword,
					}

					m.Room().AddPortal(&newExit)
					m.Room().Show(nil, fmt.Sprintf("\r\n%s falls from the sky!", newExit.Name()))
					success = true
				}
			}
		}

	case "*shutdown":
		if m.SecClearance >= c.security {
			SaveAndShutdownServer(library)
		}

	default:
		log.Printf("Command with trigger '%s' not found.", c.Trigger())
	}

	if success {
		// handle timer and usage cost
		if c.CostType() == CostFat {
			m.CurState().adjFat(-c.UseCost(), m.MaxState().fat())
		}
	}
	return success
}
