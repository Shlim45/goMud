package main

import "fmt"

type Area struct {
	Name  string
	Realm Realm
	Rooms []*Room
}

func (a *Area) SaveAreaToDBQuery() string {
	return fmt.Sprintf("INSERT INTO Area VALUES ('%s', %d)",
		a.Name, uint8(a.Realm))
}
