package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type DBConnection interface {
	DatabaseConnection() *sql.DB
	Initialize()
	Update(query string) bool
	Query(query string) *sql.Rows
	LoadAreas(w *World)
	LoadRooms(w *World)
	LoadExits(w *World)
	LoadRaces(lib *MudLib)
	LoadClasses(lib *MudLib)
	LoadAccounts(w *World)
	LoadPlayers(w *World, lib *MudLib)
	SaveRooms(w *World)
	SaveAreas(w *World)
	SavePlayers(w *World)
	SaveAccounts(w *World)
	SaveCharClasses(lib *MudLib)
	SaveRaces(lib *MudLib)
}

type DatabaseConnection struct {
	DB *sql.DB
}

func NewDatabase() *DatabaseConnection {
	return &DatabaseConnection{
		DB: DBConnect(),
	}
}

func DBConnect() *sql.DB {
	db, err := sql.Open("mysql", "mudhost:B@ckstab69@tcp(127.0.0.1:3306)/gomud")
	if err != nil {
		panic(err.Error())
	}

	return db
}

func (db *DatabaseConnection) DatabaseConnection() *sql.DB {
	return db.DB
}

func (db *DatabaseConnection) Initialize() {
	db.Update(CreateAccountTableDBQuery())

	db.Update(CreateCharClassTableDBQuery())

	db.Update(CreateRaceTableDBQuery())

	db.Update(CreateAreaTableDBQuery())

	db.Update(CreateRoomTableDBQuery())

	db.Update(CreatePortalTableDBQuery())

	db.Update(CreateItemTableDBQuery())

	db.Update(CreatePlayerTableDBQuery())
}

func (db *DatabaseConnection) Update(query string) bool {
	insert, err := db.DB.Query(query)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer insert.Close()
	return true
}

func (db *DatabaseConnection) Query(query string) *sql.Rows {
	results, err := db.DB.Query(query)
	if err != nil {
		log.Println(err.Error())
	}
	return results
}

type RaceTag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (db *DatabaseConnection) LoadRaces(lib *MudLib) {
	log.Println("Loading races...")
	results := db.Query("SELECT * FROM Race")
	count := 0
	for results.Next() {
		var raceTag RaceTag
		err := results.Scan(&raceTag.Name, &raceTag.Enabled)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		newRace := PlayerRace{
			name:        raceTag.Name,
			enabled:     raceTag.Enabled,
			statBonuses: [6]float64{},
		}
		lib.AddRace(&newRace)
		count++
	}
	defer results.Close()

	log.Printf("Done loading %d races.\r\n", count)
}

type ClassTag struct {
	Name    string `json:"name"`
	Realm   uint8  `json:"realm"`
	Enabled bool   `json:"enabled"`
}

func (db *DatabaseConnection) LoadClasses(lib *MudLib) {
	log.Println("Loading classes...")
	results := db.Query("SELECT * FROM CharClass")
	count := 0
	for results.Next() {
		var classTag ClassTag
		err := results.Scan(&classTag.Name, &classTag.Realm, &classTag.Enabled)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		newClass := PlayerClass{
			name:        classTag.Name,
			realm:       Realm(classTag.Realm),
			enabled:     classTag.Enabled,
			statBonuses: [6]float64{},
		}
		lib.AddCharClass(&newClass)
		count++
	}
	defer results.Close()

	log.Printf("Done loading %d classes.\r\n", count)
}

type AccountTag struct {
	Username string `json:"username"`
	Password string `json:"password"`
	MaxChars uint8  `json:"max_chars"`
	LastIP   string `json:"last_ip"`
	LastDate string `json:"last_date"`
	Email    string `json:"email"`
}

func (db *DatabaseConnection) LoadAccounts(w *World) {
	log.Println("Loading user accounts...")
	results := db.Query("SELECT * FROM Account")
	count := 0
	for results.Next() {
		var accTag AccountTag
		err := results.Scan(&accTag.Username, &accTag.Password, &accTag.MaxChars, &accTag.LastIP, &accTag.LastDate, &accTag.Email)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		account := NewAccount()
		account.SetUserName(accTag.Username)
		account.SetPasswordHash(accTag.Password)
		account.SetMaxChars(accTag.MaxChars)
		account.SetLastIP(accTag.LastIP)
		lastDate, err := time.Parse("2006-01-02 15:04:05.999999", accTag.LastDate)
		if err == nil {
			account.SetLastDate(lastDate)
		}
		account.SetEmail(accTag.Email)
		w.accounts[account.UserName()] = account
		count++
	}
	defer results.Close()

	log.Printf("Done loading %d user accounts.\r\n", count)
}

func (db *DatabaseConnection) LoadPlayers(w *World, lib *MudLib) {
	log.Println("Loading players...")
	results := db.Query("SELECT * FROM Player")
	count := 0
	for results.Next() {
		var pTag PlayerDB
		err := results.Scan(&pTag.name, &pTag.account, &pTag.class, &pTag.race, &pTag.room, &pTag.coins, &pTag.stre, &pTag.cons, &pTag.agil, &pTag.dext, &pTag.inte, &pTag.wisd, &pTag.con_loss, &pTag.level, &pTag.exp, &pTag.rp, &pTag.hits, &pTag.fat, &pTag.power, &pTag.trains, &pTag.guild, &pTag.guild_rank, &pTag.last_date)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var stats [NUM_STATS]uint8
		stats[STAT_STRENGTH] = pTag.stre
		stats[STAT_CONSTITUTION] = pTag.cons
		stats[STAT_AGILITY] = pTag.agil
		stats[STAT_DEXTERITY] = pTag.dext
		stats[STAT_INTELLIGENCE] = pTag.inte
		stats[STAT_WISDOM] = pTag.wisd

		baseCStats := CharStats{
			stats:     stats,
			charClass: lib.FindCharClass(pTag.class),
			race:      lib.FindRace(pTag.race),
		}

		basePStats := PhyStats{
			Attack:       uint16(3 * (stats[STAT_DEXTERITY] / 4)),
			Damage:       uint16(3 * (stats[STAT_STRENGTH] / 4)),
			Evasion:      uint16(stats[STAT_AGILITY] / 2),
			Defense:      uint16(stats[STAT_CONSTITUTION] / 2),
			MagicAttack:  uint16(stats[STAT_WISDOM]),
			MagicDamage:  uint16(stats[STAT_INTELLIGENCE]),
			MagicEvasion: uint16(3 * (stats[STAT_WISDOM] / 4)),
			MagicDefense: uint16(3 * (stats[STAT_INTELLIGENCE] / 4)),
			Level:        pTag.level,
		}

		maxCState := CharState{
			Hits:     pTag.hits,
			Fat:      pTag.fat,
			Power:    pTag.power,
			Alive:    true,
			Standing: true,
			Sitting:  false,
			Laying:   false,
		}

		newPlayer := Player{
			name:          pTag.name,
			Account:       pTag.account,
			User:          nil,
			room:          w.GetRoomById(pTag.room),
			curState:      maxCState.copyOf(),
			maxState:      &maxCState,
			basePhyStats:  &basePStats,
			curPhyStats:   basePStats.copyOf(),
			baseCharStats: &baseCStats,
			curCharStats:  baseCStats.copyOf(),
			Experience:    pTag.exp,
			RealmPoints:   pTag.rp,
			inventory:     nil,
			coins:         pTag.coins,
			tickType:      TickStop,
			tickCount:     0,
			victim:        nil,
			LastDate:      pTag.last_date,
		}
		w.characters[newPlayer.Name()] = &newPlayer
		w.accounts[pTag.account].characters[newPlayer.Name()] = &newPlayer
		count++
	}
	defer results.Close()

	log.Printf("Done loading %d players.\r\n", count)
}

type AreaTag struct {
	Name  string `json:"name"`
	Realm uint8  `json:"realm"`
}

func (db *DatabaseConnection) LoadAreas(w *World) {
	log.Println("Loading areas...")
	results := db.Query("SELECT * FROM Area")
	count := 0
	for results.Next() {
		var areaTag AreaTag
		err := results.Scan(&areaTag.Name, &areaTag.Realm)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		newArea := Area{
			Name:  areaTag.Name,
			Realm: Realm(areaTag.Realm),
		}
		w.AddArea(&newArea)
		count++
	}
	defer results.Close()

	log.Printf("Done loading %d areas.\r\n", count)
}

type RoomTag struct {
	ID    string `json:"room_id"`
	Area  string `json:"area"`
	Desc  string `json:"description"`
	Links string `json:"links"`
}

func (db *DatabaseConnection) LoadRooms(w *World) {
	log.Println("Loading rooms...")

	results := db.Query("SELECT * FROM Room")
	count := 0
	for results.Next() {
		var roomTag RoomTag
		err := results.Scan(&roomTag.ID, &roomTag.Area, &roomTag.Desc, &roomTag.Links)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var roomLinks []*RoomLink
		for n, room := range strings.Split(roomTag.Links, ";") {
			if len(room) == 0 {
				continue
			}
			link := RoomLink{
				Verb:   Direction(n).Verb(),
				RoomId: room,
			}
			roomLinks = append(roomLinks, &link)
		}

		newRoom := Room{
			Id:      roomTag.ID,
			Desc:    roomTag.Desc,
			Area:    w.areas[roomTag.Area],
			Links:   roomLinks,
			Portals: nil,
			Items:   nil,
			Mobs:    nil,
		}

		w.rooms = append(w.rooms, &newRoom)
		if area := newRoom.Area; area != nil {
			area.Rooms = append(area.Rooms, &newRoom)
		}
		count++
	}

	defer results.Close()
	log.Printf("Done loading %d rooms.\r\n", count)
}

func (db *DatabaseConnection) LoadExits(w *World) {
	log.Println("Loading exits...")

	count := 0
	for _, room := range w.rooms {
		pResults := db.Query(fmt.Sprintf("SELECT * FROM Portal WHERE room='%s'", room.RoomID()))
		for pResults.Next() {

			var portTag PortalTag
			err := pResults.Scan(&portTag.Name, &portTag.Room, &portTag.DestRoom)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			var art string
			if strings.HasPrefix(portTag.Name, "a ") {
				art = "a"
			} else if strings.HasPrefix(portTag.Name, "an ") {
				art = "an"
			} else if strings.HasPrefix(portTag.Name, "the ") {
				art = "the"
			} else if strings.HasPrefix(portTag.Name, "A ") {
				art = "A"
			} else if strings.HasPrefix(portTag.Name, "An ") {
				art = "An"
			} else if strings.HasPrefix(portTag.Name, "The ") {
				art = "The"
			} else {
				art = ""
			}

			var keyWord string
			words := strings.Split(portTag.Name, " ")
			keyWord = words[len(words)-1]

			noArtName := portTag.Name[len(art)+1:]

			newPort := Portal{
				room:     room,
				destRoom: w.GetRoomById(portTag.DestRoom),
				article:  art,
				name:     noArtName,
				keyword:  keyWord,
			}
			room.Portals = append(room.Portals, &newPort)
			count++
		}
	}

	log.Printf("Done loading %d exits.\r\n", count)
}

func (db *DatabaseConnection) SaveRooms(w *World) {
	log.Println("Saving rooms...")
	for _, room := range w.rooms {
		db.Update(room.SaveRoomToDBQuery())
		// update exits
		for _, portal := range room.Portals {
			db.Update(portal.SaveExitToDBQuery())
		}
	}
	log.Println("Done saving rooms.")
}

func (db *DatabaseConnection) SaveAreas(w *World) {
	log.Println("Saving areas...")
	for _, area := range w.areas {
		db.Update(area.SaveAreaToDBQuery())
	}
	log.Println("Done saving areas.")
}

func (db *DatabaseConnection) SavePlayers(w *World) {
	log.Println("Saving players...")
	for _, c := range w.characters {
		update, err := c.SavePlayerToDBQuery()
		if err == nil {
			db.Update(update)
		}
	}
	log.Println("Done saving players.")
}

func (db *DatabaseConnection) SaveAccounts(w *World) {
	log.Println("Saving user accounts...")
	for _, acc := range w.accounts {
		db.Update(acc.SaveAccountToDBQuery())
	}
	log.Println("Done saving accounts.")
}

func (db *DatabaseConnection) SaveCharClasses(lib *MudLib) {
	log.Println("Saving classes...")
	for _, cClass := range lib.CharClasses() {
		db.Update(cClass.SaveCharClassToDBQuery())
	}
	log.Println("Done saving classes.")
}

func (db *DatabaseConnection) SaveRaces(lib *MudLib) {
	log.Println("Saving races...")
	for _, race := range lib.Races() {
		db.Update(race.SaveRaceToDBQuery())
	}
	log.Println("Done saving races.")
}
