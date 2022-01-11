package main

import "fmt"

//
//type Item interface {
//	Owner() *ItemPossessor
//	SetOwner(newOwner *ItemPossessor)
//}
//
//type ItemPossessor interface {
//	AddItem(item *MudItem)
//	RemoveItem(item *MudItem)
//	MoveItemTo(item *MudItem)
//}

type Item struct {
	keyword  string
	article  string
	name     string
	owner    *Player
	location *Room
}

func (item *Item) Owner() *Player {
	return item.owner
}

func (item *Item) SetOwner(newOwner *Player) {
	item.owner = newOwner
}

func (item *Item) Location() *Room {
	return item.location
}

func (item *Item) SetLocation(newLocation *Room) {
	item.location = newLocation
}

func (item *Item) FullName() string {
	if len(item.article) > 0 {
		return fmt.Sprintf("%s %s", item.article, item.name)
	}
	return item.name
}

func (item *Item) Name() string {
	return item.name
}

func (item *Item) SetName(newName string) {
	item.name = newName
}

func (item *Item) Article() string {
	return item.article
}

func (item *Item) SetArticle(newArt string) {
	item.article = newArt
}

func (item *Item) Keyword() string {
	return item.keyword
}

func (item *Item) SetKeyword(newKey string) {
	item.keyword = newKey
}
