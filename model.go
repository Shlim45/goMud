package main

type User struct {
	Session   *Session
	Account   *Account
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
	REALM_IMMORTAL = iota
	REALM_EVIL
	REALM_CHAOS
	REALM_GOOD
	REALM_KAID
)

func (R Realm) String() string {
	switch R {
	case REALM_IMMORTAL:
		return "Immortal"
	case REALM_EVIL:
		return "Evil"
	case REALM_CHAOS:
		return "Chaos"
	case REALM_GOOD:
		return "Good"
	case REALM_KAID:
		return "Kaid"
	default:
		return "None"
	}
}

func (R Realm) God() string {
	switch R {
	case REALM_IMMORTAL:
		return "Xyz"
	case REALM_EVIL:
		return "Arnak"
	case REALM_CHAOS:
		return "Ra'Kur"
	case REALM_GOOD:
		return "Niord"
	case REALM_KAID:
		return "Abc"
	default:
		return "None"
	}
}

type MudMob interface {
	Name() string
	SetName(newName string)
	basePhyStats() *PhyStats
	curPhyStats() *PhyStats
	recoverPhyStats()
	curState() *CharState
	maxState() *CharState
	adjHits(amount uint16)
	adjFat(amount uint16)
	adjPower(amount uint16)
	adjMaxHits(amount uint16)
	adjMaxFat(amount uint16)
	adjMaxPower(amount uint16)
	recoverCharState()
	curCharStats() *CharStats
	baseCharStats() *CharStats
	recoverCharStats()
	level()
	setLevel(newLevel uint8)
	isPlayer() bool
	Inventory() []*Item
	AddItem(item *Item)
	RemoveItem(item *Item)
	MoveItemTo(item *Item)
	AttackVictim()
	Tick(tType TickType) bool
	Init(library *MudLib)
	Walk(dest *Room, verb string)
	WalkThrough(port *Portal)
	attackTarget(target *MOB)
	killMOB(killer *MOB)
	damageMOB(attacker *MOB, dmg uint16)
}
