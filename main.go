package main

import (
	"log"
)

func main() {
	ch := make(chan SessionEvent)

	db := NewDatabase()
	db.Initialize()

	world := NewWorld(db)
	//world.Init()

	db.LoadAreas(world)
	db.LoadRooms(world)

	library := NewLibrary(world)
	library.LoadCommands()
	library.LoadRaces()
	library.LoadCharClasses()

	sessionHandler := NewSessionHandler(world, library, ch)
	go sessionHandler.Start()

	err := startServer(ch)
	if err != nil {
		log.Fatalln(err)
	}
}
