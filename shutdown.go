package main

import "log"

func SaveAndShutdownServer(world *World, library *MudLib) {
	// TODO(jon): This does not UPDATE, only INSERT.  Need to save changes.
	log.Println("Shutting down server...")
	world.db.SaveRaces(library)
	world.db.SaveCharClasses(library)
	world.db.SaveAccounts(world)
	world.db.SavePlayers(world)
	world.db.SaveAreas(world)
	world.db.SaveRooms(world)
	log.Println("World and objects saved.")

	log.Println("\r\nLogging off players.")
	for _, mob := range world.characters {
		if mob.User.Session.Status() == INGAME {
			world.RemoveFromWorld(mob)
		}
		mob.User.Session.WriteLine("\r\nGoodbye")
		mob.User.Session.conn.Close()
	}
	log.Println("All players logged off.")

	log.Println("\r\nClosing database connection...")
	defer world.db.DatabaseConnection().Close()
	log.Println("Database disconnected.")

	// TODO(jon): Need a flag on server? that can be set to shut game down.
}
