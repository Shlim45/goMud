package main

type User struct {
	Session   *Session
	Account   *Account
	Character *Player
	ANSI      bool
}
