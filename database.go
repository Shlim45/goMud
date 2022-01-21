package main

import (
	"database/sql"
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
	LoadRooms(w *World)
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

	// TODO(jon): handle during shutdown
	//defer db.Close()
	return db
}

func (db *DatabaseConnection) DatabaseConnection() *sql.DB {
	return db.DB
}

func (db *DatabaseConnection) Initialize() {
	db.Update("CREATE TABLE IF NOT EXISTS Account(" +
		"username VARCHAR(20) PRIMARY KEY," +
		"password CHAR(60) NOT NULL," +
		"max_chars TINYINT UNSIGNED NOT NULL DEFAULT 3," +
		"last_ip VARCHAR(15)," +
		"last_date TIMESTAMP," +
		"email VARCHAR(319) NOT NULL" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS CharClass(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"realm TINYINT NOT NULL DEFAULT 0," +
		"enabled BOOLEAN NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Race(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"enabled BOOLEAN NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Area(" +
		"name VARCHAR(50) PRIMARY KEY," +
		"realm TINYINT UNSIGNED NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Room(" +
		"room_id VARCHAR(5) NOT NULL," +
		"area VARCHAR(50) NOT NULL," +
		"description VARCHAR(512)," +
		"links VARCHAR(1024)," +
		"PRIMARY KEY (area, room_id)," +
		"FOREIGN KEY (area) REFERENCES Area(name) ON UPDATE CASCADE ON DELETE RESTRICT" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Portal(" +
		"name VARCHAR(50) NOT NULL," +
		"room VARCHAR(60) NOT NULL," +
		"dest_room VARCHAR(60)," +
		"PRIMARY KEY (name, room)" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Item(" +
		"article VARCHAR(3)," +
		"name VARCHAR(50) NOT NULL," +
		"keyword VARCHAR(50) NOT NULL," +
		"owner VARCHAR(60) NOT NULL," +
		"value INT UNSIGNED NOT NULL," +
		"item_type TINYINT UNSIGNED NOT NULL," +
		"PRIMARY KEY (article, name)" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Player(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"account VARCHAR(20)," +
		"class VARCHAR(20)," +
		"race VARCHAR(20)," +
		"room VARCHAR(60)," +
		"coins BIGINT UNSIGNED NOT NULL," +
		"stre TINYINT UNSIGNED NOT NULL," +
		"cons TINYINT UNSIGNED NOT NULL," +
		"agil TINYINT UNSIGNED NOT NULL," +
		"dext TINYINT UNSIGNED NOT NULL," +
		"inte TINYINT UNSIGNED NOT NULL," +
		"wisd TINYINT UNSIGNED NOT NULL," +
		"con_loss TINYINT UNSIGNED NOT NULL," +
		"level TINYINT UNSIGNED NOT NULL," +
		"exp BIGINT UNSIGNED NOT NULL," +
		"rp INT UNSIGNED NOT NULL," +
		"hits SMALLINT UNSIGNED NOT NULL," +
		"fat SMALLINT UNSIGNED NOT NULL," +
		"power SMALLINT UNSIGNED NOT NULL," +
		"trains SMALLINT UNSIGNED NOT NULL DEFAULT 0," +
		"guild VARCHAR(30)," + // FK
		"guild_rank TINYINT UNSIGNED NOT NULL DEFAULT 0," +
		"last_date TIMESTAMP," +
		"FOREIGN KEY (account) REFERENCES Account(username) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (class) REFERENCES CharClass(name) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (race) REFERENCES Race(name) ON UPDATE CASCADE ON DELETE SET NULL" +
		")")
}

func (db *DatabaseConnection) Update(query string) bool {
	insert, err := db.DB.Query(query)
	if err != nil {
		log.Println(err.Error())
		//panic(err.Error())
		return false
	}
	defer insert.Close()
	return true
}

func (db *DatabaseConnection) Query(query string) *sql.Rows {
	results, err := db.DB.Query(query)
	if err != nil {
		log.Println(err.Error())
		//panic(err.Error())
	}
	return results
}

type AreaTag struct {
	Name  string `json:"name"`
	Realm uint8  `json:"realm"`
}

func (db *DatabaseConnection) LoadAreas(w *World) {
	log.Println("Loading areas...")
	results := db.Query("SELECT * FROM Area")

	for results.Next() {
		var areaTag AreaTag
		err := results.Scan(&areaTag.Name, &areaTag.Realm)
		if err != nil {
			log.Println(err.Error())
			//panic(err.Error())
		}

		newArea := Area{
			Name:  areaTag.Name,
			Realm: Realm(areaTag.Realm),
		}
		w.AddArea(&newArea)
	}
	defer results.Close()

	log.Println("Done loading areas.")
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

	for results.Next() {
		var roomTag RoomTag
		err := results.Scan(&roomTag.ID, &roomTag.Area, &roomTag.Desc, &roomTag.Links)
		if err != nil {
			log.Println(err.Error())
			//panic(err.Error())
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
		/*
			for n, r := range strings.Split(roomTag.Links, ",") {
				if len(r) > 0 {
					roomLink := RoomLink{
						Verb:   Direction(n).Verb(),
						RoomId: r,
					}
					newRoom.Links = append(newRoom.Links, &roomLink)
				}
			}
		*/
		w.rooms = append(w.rooms, &newRoom)
	}

	defer results.Close()
	log.Println("Done loading rooms...")
}

func (db *DatabaseConnection) SaveRooms(w *World) {
	log.Println("Saving rooms...")
	for _, room := range w.rooms {
		db.Update(room.SaveRoomToDBQuery())
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

// TimeString - given a time, return the MySQL standard string representation
func TimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999999")
}
