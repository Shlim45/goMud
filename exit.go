package main

import (
	"fmt"
	"strings"
)

type Portal struct {
	room     *Room `json:"room"`
	destRoom *Room `json:"dest_room"`
	article  string
	name     string
	keyword  string
}

func (p *Portal) Location() *Room {
	return p.room
}

func (p *Portal) SetLocation(newRoom *Room) {
	p.room = newRoom
}

func (p *Portal) DestRoom() *Room {
	return p.destRoom
}

func (p *Portal) SetDestRoom(newDest *Room) {
	p.destRoom = newDest
}

func (p *Portal) Article() string {
	return p.article
}

func (p *Portal) SetArticle(newArt string) {
	p.article = newArt
}

func (p *Portal) Name() string {
	return p.name
}

func (p *Portal) SetName(newName string) {
	p.name = newName
}

func (p *Portal) Keyword() string {
	return p.keyword
}

func (p *Portal) SetKeyword(newKey string) {
	p.keyword = newKey
}

func (p *Portal) FullName() string {
	var full strings.Builder
	if len(p.article) > 0 {
		full.WriteString(p.article + " ")
	}
	full.WriteString(p.name)
	return full.String()
}

type PortalTag struct {
	Name     string `json:"name"`
	Room     string `json:"room"`
	DestRoom string `json:"dest_room"`
}

/**
describe Portal;
+-----------+-------------+------+-----+---------+-------+
| Field     | Type        | Null | Key | Default | Extra |
+-----------+-------------+------+-----+---------+-------+
| name      | varchar(50) | NO   | PRI | NULL    |       |
| room      | varchar(60) | NO   | PRI | NULL    |       |
| dest_room | varchar(60) | YES  |     | NULL    |       |
+-----------+-------------+------+-----+---------+-------+
*/

func (p *Portal) SaveExitToDBQuery() string {
	portTag := PortalTag{
		Name:     p.FullName(),
		Room:     p.Location().RoomID(),
		DestRoom: p.DestRoom().RoomID(),
	}

	return fmt.Sprintf("INSERT INTO Portal VALUES ('%s', '%s', '%s') AS new "+
		"ON DUPLICATE KEY UPDATE name=new.name, room=new.room, dest_room=new.dest_room",
		portTag.Name, portTag.Room, portTag.DestRoom)
}

func CreatePortalTableDBQuery() string {
	return "CREATE TABLE IF NOT EXISTS Portal(" +
		"name VARCHAR(50) NOT NULL," +
		"room VARCHAR(60) NOT NULL," +
		"dest_room VARCHAR(60)," +
		"PRIMARY KEY (name, room)" +
		")"
}
