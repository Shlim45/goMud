package main

import "fmt"

type MudItem interface {
	Owner() ItemPossessor
	SetOwner(newOwner ItemPossessor)
	Article() string
	SetArticle(newArt string)
	Name() string
	SetName(newName string)
	Keyword() string
	SetKeyword(newKey string)
	FullName() string
}

type ItemPossessor interface {
	AddItem(item *Item)
	RemoveItem(item *Item)
	MoveItemTo(item *Item)
}

type ItemType uint8

const (
	TYPE_GENERIC = iota
	TYPE_WEAPON
	TYPE_ARMOR
	TYPE_COINS
)

type Item struct {
	keyword  string
	article  string
	name     string
	owner    ItemPossessor
	value    uint64
	itemType ItemType
}

func (item *Item) Value() uint64 {
	return item.value
}

func (item *Item) SetValue(newValue uint64) {
	item.value = newValue
}

func (item *Item) ItemType() ItemType {
	return item.itemType
}

func (item *Item) SetItemType(newType ItemType) {
	item.itemType = newType
}

func (item *Item) Owner() ItemPossessor {
	return item.owner
}

func (item *Item) SetOwner(newOwner ItemPossessor) {
	item.owner = newOwner
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

/*
Portals: []*RoomLink{
		{
			Verb:   "gate",
			RoomId: "D",
		},
	}
*/

func CreateItemTableDBQuery() string {
	return "CREATE TABLE IF NOT EXISTS Item(" +
		"article VARCHAR(3)," +
		"name VARCHAR(50) NOT NULL," +
		"keyword VARCHAR(50) NOT NULL," +
		"owner VARCHAR(60) NOT NULL," +
		"value INT UNSIGNED NOT NULL," +
		"item_type TINYINT UNSIGNED NOT NULL," +
		"PRIMARY KEY (article, name)" +
		")"
}
