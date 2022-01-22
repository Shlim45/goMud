package main

import "fmt"

type Area struct {
	Name  string
	Realm Realm
	Rooms []*Room
}

func (a *Area) SaveAreaToDBQuery() string {
	return fmt.Sprintf("INSERT INTO Area VALUES ('%s', %d) AS new ON DUPLICATE KEY UPDATE name=new.name, realm=new.realm",
		a.Name, uint8(a.Realm))
}

func (a *Area) GetRoomById(id string) *Room {
	for _, r := range a.Rooms {
		if r.Id == id {
			return r
		}
	}
	return nil
}
