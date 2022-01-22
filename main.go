package main

import (
	"log"
)

func main() {
	ch := make(chan SessionEvent)

	db := NewDatabase()
	db.Initialize()

	world := NewWorld(db)

	db.LoadAreas(world)
	db.LoadRooms(world)
	db.LoadExits(world)

	library := NewLibrary(world)
	library.LoadCommands()
	db.LoadRaces(library)
	db.LoadClasses(library)
	db.LoadAccounts(world)
	db.LoadPlayers(world, library)

	world.Tick()

	sessionHandler := NewSessionHandler(world, library, ch)
	go sessionHandler.Start()

	err := startServer(ch)
	if err != nil {
		log.Fatalln(err)
	}
}
