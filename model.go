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

type TimerType uint8

const (
	TIMER_NONE = iota
	TIMER_ATTACK
	TIMER_SPELL
	TIMER_SKILL
)

type ActionCost uint8

const (
	COST_NONE = iota
	COST_HITS
	COST_FAT
	COST_POWER
)

type Realm uint8

const (
	IMMORTAL = iota
	EVIL
	CHAOS
	GOOD
	KAID
)

func (R Realm) String() string {
	switch R {
	case IMMORTAL:
		return "Immortal"
	case EVIL:
		return "Evil"
	case CHAOS:
		return "Chaos"
	case GOOD:
		return "Good"
	case KAID:
		return "Kaid"
	default:
		return "None"
	}
}

func (R Realm) God() string {
	switch R {
	case IMMORTAL:
		return "Xyz"
	case EVIL:
		return "Arnak"
	case CHAOS:
		return "Ra'Kur"
	case GOOD:
		return "Niord"
	case KAID:
		return "Abc"
	default:
		return "None"
	}
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
