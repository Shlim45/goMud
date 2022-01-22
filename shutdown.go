package main

import "log"

func SaveAndShutdownServer(world *World, library *MudLib) {
	world.Broadcast(CHighlight("Shutdown initiated."))
	log.Println("Shutting down server...")
	world.db.SaveRaces(library)
	world.db.SaveCharClasses(library)
	world.db.SaveAccounts(world)
	world.db.SavePlayers(world)
	world.db.SaveAreas(world)
	world.db.SaveRooms(world)
	log.Println("World and objects saved.")

	log.Println("Logging off players.")
	for _, mob := range world.characters {
		if u := mob.User; u != nil {
			if u.Session.Status() == INGAME {
				world.RemoveFromWorld(mob)
			}
			u.Session.WriteLine("\r\nGoodbye")
			u.Session.conn.Close()
		}
	}
	log.Println("All players logged off.")

	log.Println("Closing database connection...")
	defer world.db.DatabaseConnection().Close()
	log.Println("Database disconnected.")

	// TODO(jon): Need a flag on server? that can be set to shut game down.
}
