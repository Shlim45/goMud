package main

import "log"

func SaveAndShutdownServer(world *World, library *MudLib) {
	log.Println("Shutting down server...")
	world.db.SaveRaces(library)
	world.db.SaveCharClasses(library)
	world.db.SaveAccounts(world)
	world.db.SavePlayers(world)
	world.db.SaveAreas(world)
	world.db.SaveRooms(world)
	log.Println("World and objects saved.")

	log.Println("Closing database connection...")
	defer world.db.DatabaseConnection().Close()
	log.Println("Database disconnected.")
}
