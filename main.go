package main

import (
	"log"
)

func main() {
	ch := make(chan SessionEvent)

	db := NewDatabase()
	db.Initialize()

	world := NewWorld()
	world.Init()

	library := NewLibrary(world)
	library.LoadCommands()
	library.LoadCharClasses()

	sessionHandler := NewSessionHandler(world, library, ch)
	go sessionHandler.Start()

	err := startServer(ch)
	if err != nil {
		log.Fatalln(err)
	}
}
