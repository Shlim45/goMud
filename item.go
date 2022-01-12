package main

import "fmt"

type ItemPossessor interface {
	AddItem(item *Item)
	RemoveItem(item *Item)
	MoveItemTo(item *Item)
}

type Item struct {
	keyword string
	article string
	name    string
	owner   *ItemPossessor
}

func (item *Item) Owner() *ItemPossessor {
	return item.owner
}

func (item *Item) SetOwner(newOwner ItemPossessor) {
	item.owner = &newOwner
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
