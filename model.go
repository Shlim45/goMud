package main

type Entity struct {
	entityId string
}

func (e *Entity) EntityId() string {
	return e.entityId
}

type User struct {
	Session   *Session
	Character *Player
	ANSI      bool
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

type CharState struct {
	Hits     uint16
	Fat      uint16
	Power    uint16
	Alive    bool
	Standing bool
	Sitting  bool
	Laying   bool
}

func (cState *CharState) copyOf() *CharState {
	copyOf := CharState{
		Hits:     cState.Hits,
		Fat:      cState.Fat,
		Power:    cState.Power,
		Alive:    cState.Alive,
		Standing: cState.Standing,
		Sitting:  cState.Sitting,
		Laying:   cState.Laying,
	}
	return &copyOf
}
