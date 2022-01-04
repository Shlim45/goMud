package main

type RoomLink struct {
	Verb   string
	RoomId string
}

type Room struct {
	Id    string
	Desc  string
	Links []*RoomLink

	Characters []*Character
}

func (r *Room) Show(source *Character, msg string) {
	for _, player := range r.Characters {
		player.SendMessage(msg)
	}
}

func (r *Room) ShowOthers(source *Character, msg string) {
	for _, player := range r.Characters {
		if player != nil && player != source {
			player.SendMessage(msg)
		}
	}
}

func (r *Room) AddCharacter(character *Character) {
	r.Characters = append(r.Characters, character)
	character.Room = r
}

func (r *Room) RemoveCharacter(character *Character) {
	character.Room = nil

	var characters []*Character
	for _, c := range r.Characters {
		if c != character {
			characters = append(characters, c)
		}
	}
	r.Characters = characters
}
