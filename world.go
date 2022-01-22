package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type World struct {
	accounts    map[string]*Account
	charClasses map[string]CharClass
	characters  map[string]*Player
	areas       map[string]*Area
	rooms       []*Room
	db          DBConnection
	tickCount   uint64
}

func NewWorld(db DBConnection) *World {
	return &World{
		accounts:    make(map[string]*Account),
		charClasses: make(map[string]CharClass),
		characters:  make(map[string]*Player),
		areas:       make(map[string]*Area),
		db:          db,
	}
}

func (w *World) AddArea(area *Area) {
	_, exists := w.areas[area.Name]
	if !exists {
		w.areas[area.Name] = area
	}
}

func (w *World) Tick() {
	if w.tickCount > 0 {
		if w.tickCount%(60*5) == 0 {
			w.db.SavePlayers(w)
		}

	}
	w.tickCount++
	time.Sleep(1000 * time.Millisecond)
	go w.Tick()
}

func (w *World) HandleCharacterJoined(character *Player) {
	room := character.Room()
	if room != nil {
		room.AddMOB(character)
	} else {
		w.rooms[0].AddMOB(character)
	}

	character.Room().ShowRoom(character)
	character.Room().ShowOthers(character, nil, fmt.Sprintf("%s appears in a puff of smoke.", character.Name()))

	log.Println(fmt.Sprintf("Player login: %s", character.Name()))
}

func (w *World) RemoveFromWorld(character *Player) {
	room := character.Room()
	if room != nil {
		room.RemoveMOB(character)
		room.Show(nil, fmt.Sprintf("%s disappears in a puff of smoke.", character.Name()))
	}
	character.LastDate = TimeString(time.Now())
	if account := w.accounts[character.Account]; account != nil {
		account.UpdateLastDate()
	}
	log.Println(fmt.Sprintf("Player logout: %s", character.Name()))
}

func (w *World) Broadcast(msg string) {
	for _, player := range w.characters {
		if player.User != nil && player.User.Session.Status() == INGAME {
			player.SendMessage(msg, true)
		}
	}
}

func (w *World) GetRoomById(id string) *Room {
	if strings.Contains(id, "#") {
		areaName := before(id, "#")
		roomId := after(id, "#")
		area := w.areas[areaName]
		if area != nil {
			return area.GetRoomById(roomId)
		}
	}
	for _, r := range w.rooms {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func (w *World) MoveMob(mob MOB, to *Room) {
	mob.Room().RemoveMOB(mob)
	to.AddMOB(mob)
	to.ShowRoom(mob)
	mob.adjFat(-2, mob.maxState().fat())
}

func (w *World) FindUserAccount(uName string) (*Account, error) {
	uName = strings.ToLower(uName)
	for key, account := range w.accounts {
		lowKey := strings.ToLower(key)
		if strings.Compare(lowKey, uName) == 0 {
			return account, nil
		}
	}
	return nil, errors.New("user Account not found")
}

func (w *World) HandleUserLogin(user *User, input string) {
	session := user.Session
	switch session.Status() {
	case DEFAULT:
		session.WriteLine(CHighlight("\r\nWelcome to Darkness Falls.\r\n"))
		session.WriteLine("Username: ")
		session.SetStatus(USERNAME)

	case USERNAME:
		account, err := w.FindUserAccount(input)
		if err == nil {
			user.Account = account
		}
		session.WriteLine("\r\nPassword: ")
		// TODO(jon): TELNET ECHOOFF
		session.SetStatus(PASSWORD)

	case PASSWORD:
		if user.Account != nil {
			if user.Account.CheckPasswordHash(input) {
				session.SetStatus(MENU)
			} else {
				session.WriteLine("\r\nInvalid username and/or password.")
				user.Account = nil
				session.WriteLine("\r\nUsername: ")
				session.SetStatus(USERNAME)
			}
		} else {
			session.WriteLine("\r\nInvalid username and/or password.")
			session.WriteLine("\r\nUsername: ")
			session.SetStatus(USERNAME)
		}
		if session.Status() == MENU {
			session.WriteLine("\r\n")

			session.WriteLine("Please choose an option:")
			session.WriteLine(fmt.Sprintf("%s - Create Character", CHighlight("C")))
			session.WriteLine(fmt.Sprintf("%s - Select Character", CHighlight("S")))
			session.WriteLine(fmt.Sprintf("%s - Quit Darkness Falls", CHighlight("Q")))
		}

	case MENU:
		switch strings.ToUpper(input) {
		case "C":
			session.WriteLine("Coming soon!  Bye now!")
			session.SetStatus(CREATE)

		case "S":
			var chars strings.Builder
			if len(user.Account.characters) > 0 {
				chars.WriteString("\r\n")
				for _, mob := range user.Account.characters {
					chars.WriteString(fmt.Sprintf("%-20v lvl %d %s\r\n",
						BrightBlue(mob.Name()), mob.curPhyStats().level(), mob.curCharStats().CurrentClass().Name()))
				}
				chars.WriteString("\r\nEnter Character Name: ")
				session.WriteLine(chars.String())
				session.SetStatus(SELECT)
			} else {
				chars.WriteString("\r\nNo Characters found.  CREATE a new Character!")
				session.WriteLine("Coming soon!  Bye now!")
				session.SetStatus(CREATE)
			}

		case "Q":
			session.SetStatus(QUIT)

		default:
			session.WriteLine("\r\n")

			session.WriteLine(CHighlight("Please choose an option:"))
			session.WriteLine(fmt.Sprintf("%s - Create Character", CHighlight("C")))
			session.WriteLine(fmt.Sprintf("%s - Select Character", CHighlight("S")))
			session.WriteLine(fmt.Sprintf("%s - Quit Darkness Falls", CHighlight("Q")))

			session.WriteLine("\r\n")
		}

	case CREATE:
		session.SetStatus(QUIT)

	case SELECT:
		var chars strings.Builder
		for name, mob := range user.Account.characters {
			chars.WriteString(fmt.Sprintf("%-15v lvl %d %s\r\n",
				CHighlight(mob.Name()), mob.curPhyStats().level(), mob.curCharStats().CurrentClass().Name()))
			if strings.Compare(strings.ToLower(input), strings.ToLower(name)) == 0 {
				user.Character = mob
				mob.User = user
				session.SetStatus(INGAME)
				w.HandleCharacterJoined(mob)
				return
			}
		}
		chars.WriteString("\r\nInvalid selections, Enter Character Name: ")

	case QUIT:
		if strings.Compare(strings.ToLower(input), "y") == 0 {
			session.SetStatus(DEFAULT)
			session.WriteLine("\r\nGoodbye")
			if user.Character != nil {
				w.RemoveFromWorld(user.Character)
			}
			session.conn.Close()
		} else if user.Character != nil {
			session.SetStatus(INGAME)
		} else {
			session.WriteLine("\r\n")

			session.WriteLine(CHighlight("Please choose an option:"))
			session.WriteLine(fmt.Sprintf("%s - Create Character", CHighlight("C")))
			session.WriteLine(fmt.Sprintf("%s - Select Character", CHighlight("S")))
			session.WriteLine(fmt.Sprintf("%s - Quit Darkness Falls", CHighlight("Q")))

			session.WriteLine("\r\n")
			session.SetStatus(MENU)
		}
	}
}

func (w *World) HandlePlayerInput(player *Player, input string, library *MudLib) {
	tokens := strings.Split(input, " ")
	success := false
	cmd := library.FindCommand(tokens[0])

	for _, link := range player.Room().Links {
		if link.Verb == tokens[0] || strings.HasPrefix(link.Verb, tokens[0]) {
			if player.curState().fat() < 2 {
				player.SendMessage("You are too tired!", true)
				return
			}

			target := w.GetRoomById(link.RoomId)
			if target != nil {
				player.Walk(target, link.Verb)
				success = true
				break
			}
		}
	}

	if !success {
		if cmd != nil {
			go cmd.ExecuteCmd(player, tokens, w, library)
		} else {
			player.SendMessage(fmt.Sprintf("Huh?  Command '%s' not found.", tokens[0]), true)
		}
	}
}
