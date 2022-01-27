package main

import "log"

func SaveAndShutdownServer(library *MudLib) {
	library.world.Broadcast(CHighlight("Shutdown initiated."))
	log.Println("Shutting down server...")
	library.world.db.SaveRaces(library)
	library.world.db.SaveCharClasses(library)
	library.world.db.SaveAccounts(library.world)
	library.world.db.SavePlayers(library.world)
	library.world.db.SaveAreas(library.world)
	library.world.db.SaveRooms(library.world)
	log.Println("World and objects saved.")

	log.Println("Closing database connection...")
	defer library.world.db.DatabaseConnection().Close()

	log.Println("Logging off players.")
	for _, mob := range library.world.characters {
		if u := mob.User; u != nil {
			if u.Session.Status() == INGAME {
				library.world.RemoveFromWorld(mob)
			}
			u.Session.WriteLine("\r\nGoodbye")
			u.Session.conn.Close()
		}
	}
	log.Println("All players logged off.")

	// TODO(jon): Need a flag on server? that can be set to shut game down.
}
