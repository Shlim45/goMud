package main

type TickType uint8

const (
	TickNormal = iota
	TickMonster
	TickPlayer
	TickStop
)

type TimerType uint8

const (
	TimerNone = iota
	TimerAttack
	TimerSpell
	TimerSkill
)

type ActionCost uint8

const (
	CostNone = iota
	CostHits
	CostFat
	CostPower
)

type MOB interface {
	Name() string
	SetName(newName string)
	Coins() uint64
	SetCoins(amount uint64)
	AdjCoins(amount int64)
	basePhyStats() *PhyStats
	curPhyStats() *PhyStats
	recoverPhyStats()
	curState() *CharState
	maxState() *CharState
	adjHits(amount int16, max uint16)
	adjFat(amount int16, max uint16)
	adjPower(amount int16, max uint16)
	adjMaxHits(amount int16)
	adjMaxFat(amount int16)
	adjMaxPower(amount int16)
	recoverCharState()
	curCharStats() *CharStats
	baseCharStats() *CharStats
	recoverCharStats()
	isPlayer() bool
	Inventory() []*Item
	AddItem(item *Item)
	RemoveItem(item *Item)
	MoveItemTo(item *Item)
	AttackVictim()
	Tick(tType TickType) bool
	Init(library *MudLib)
	SendMessage(msg string, newLine bool)
	Walk(dest *Room, verb string)
	WalkThrough(port *Portal)
	Room() *Room
	SetRoom(newRoom *Room)
	attackTarget(target MOB)
	killMOB(killer MOB)
	damageMOB(attacker MOB, dmg uint16)
	Victim() MOB
	SetVictim(newVictim MOB)
	AwardExp(howMuch uint64)
	AwardRP(howMuch uint32)
}
