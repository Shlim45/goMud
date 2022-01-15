package main

type Entity struct {
	entityId string
}

func (e *Entity) EntityId() string {
	return e.entityId
}

type User struct {
	Session   *Session
	Character *MOB
	ANSI      bool
}

type Environmental interface {
	Name() string
	SetName(newName string)
}

type Mob interface {
	basePhyStats() *PhyStats
	curPhyStats() *PhyStats
	recoverPhyStats()
	curState() *CharState
	adjHits(amount uint16)
	adjFat(amount uint16)
	adjPower(amount uint16)
	maxState() *CharState
	adjMaxHits(amount uint16)
	adjMaxFat(amount uint16)
	adjMaxPower(amount uint16)
	recoverCharState()
	level()
	setLevel(newLevel uint8)
}
