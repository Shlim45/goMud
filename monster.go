package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Monster struct {
	name          string `json:"name"`
	room          *Room
	curState      *CharState
	maxState      *CharState
	basePhyStats  *PhyStats
	curPhyStats   *PhyStats
	baseCharStats *CharStats
	curCharStats  *CharStats
	Experience    uint64 `json:"exp"`
	inventory     []*Item
	coins         uint64 `json:"coins"`
	tickType      TickType
	tickCount     uint64
	victim        MOB
}

func (m *Monster) IsPlayer() bool {
	return false
}

func (m *Monster) Inventory() []*Item {
	return m.inventory
}

func (m *Monster) AddItem(item *Item) {
	m.inventory = append(m.inventory, item)
	item.SetOwner(m)
}

func (m *Monster) RemoveItem(item *Item) {
	item.SetOwner(nil)

	var items []*Item
	for _, i := range m.inventory {
		if i != item {
			items = append(items, i)
		}
	}
	m.inventory = items
}

func (m *Monster) MoveItemTo(item *Item) {
	if item.Owner() != nil {
		item.Owner().RemoveItem(item)
		item.SetOwner(nil)
	}

	m.inventory = append(m.inventory, item)
	item.SetOwner(m)
}

func (m *Monster) Name() string {
	return m.name
}

func (m *Monster) SetName(newName string) {
	m.name = newName
}

func (m *Monster) Coins() uint64 {
	return m.coins
}

func (m *Monster) SetCoins(amount uint64) {
	m.coins = amount
}

func (m *Monster) AdjCoins(amount int64) {
	if amount < 0 {
		less := uint64(-amount)
		if less > m.coins {
			m.coins = 0
		} else {
			m.coins -= less
		}
	} else {
		m.coins += uint64(amount)
	}
}

func (m *Monster) AttackVictim() {
	victim := m.Victim()
	if victim == nil {
		return
	}
	if victim.Room() == m.Room() {
		m.AttackTarget(victim)
	}
}

func (m *Monster) Tick(tType TickType) bool {
	if m == nil || tType == TickStop {
		return false
	}

	m.tickCount++
	switch m.tickType {
	case TickNormal:

	case TickMonster:
		if !m.CurState().alive() {
			return false
		}

		if m.tickCount%4 == 0 {
			m.AttackVictim()
		}

	case TickStop:
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	go m.Tick(m.tickType)
	return m.tickType != TickStop
}

func (m *Monster) Init(library *MudLib) {
	rand.Seed(time.Now().UnixMilli())

	var stats [NUM_STATS]uint8

	var newRace Race
	var newClass CharClass
	for i := range stats {
		stats[i] = uint8(rand.Intn(10) + 5)
	}
	newRace = library.FindRace("Human")
	newClass = library.FindCharClass("Fighter")

	baseCStats := CharStats{
		stats:     stats,
		charClass: newClass,
		race:      newRace,
	}
	m.baseCharStats = &baseCStats
	m.curCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(3 * (m.BaseCharStats().Stats()[STAT_DEXTERITY] / 4)),
		Damage:       uint16(3 * (m.BaseCharStats().Stats()[STAT_STRENGTH] / 4)),
		Evasion:      uint16(m.BaseCharStats().Stats()[STAT_AGILITY] / 2),
		Defense:      uint16(m.BaseCharStats().Stats()[STAT_CONSTITUTION] / 2),
		MagicAttack:  uint16(m.BaseCharStats().Stats()[STAT_WISDOM]),
		MagicDamage:  uint16(m.BaseCharStats().Stats()[STAT_INTELLIGENCE]),
		MagicEvasion: uint16(3 * (m.BaseCharStats().Stats()[STAT_WISDOM] / 4)),
		MagicDefense: uint16(3 * (m.BaseCharStats().Stats()[STAT_INTELLIGENCE] / 4)),
		Level:        1 + uint8(m.Experience/1000),
	}
	m.basePhyStats = &baseStats
	m.curPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16((baseStats.Level * 30) + (m.BaseCharStats().Stats()[STAT_CONSTITUTION] / 10.0)),
		Fat:      uint16((baseStats.Level * 30) + (m.BaseCharStats().Stats()[STAT_CONSTITUTION] / 10.0)),
		Power:    uint16((baseStats.Level * 20) + (m.BaseCharStats().Stats()[STAT_INTELLIGENCE] / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	m.maxState = &baseState
	m.curState = baseState.copyOf()

	// should only be stop on new mobs
	if m.tickType == TickStop {
		m.tickType = TickMonster

		go m.Tick(m.tickType)
	}
}

func (m *Monster) BasePhyStats() *PhyStats {
	return m.basePhyStats
}

func (m *Monster) CurPhyStats() *PhyStats {
	return m.curPhyStats
}

func (m *Monster) RecoverPhyStats() {
	m.curPhyStats = m.BasePhyStats().copyOf()
}

func (m *Monster) CurState() *CharState {
	return m.curState
}

func (m *Monster) MaxState() *CharState {
	return m.maxState
}

func (m *Monster) RecoverCharState() {
	m.curState = m.MaxState().copyOf()
}

func (m *Monster) CurCharStats() *CharStats {
	return m.curCharStats
}

func (m *Monster) BaseCharStats() *CharStats {
	return m.baseCharStats
}

func (m *Monster) RecoverCharStats() {
	m.curCharStats = m.BaseCharStats().copyOf()
}

func (m *Monster) SendMessage(msg string, newLine bool) {
	// TODO(jon): if possessed
}

func (m *Monster) Walk(dest *Room, verb string) {
	m.Room().ShowOthers(m, nil, fmt.Sprintf("%s went %s.", m.Name(), verb))
	m.Room().RemoveMOB(m)
	dest.AddMOB(m)
	dest.ShowRoom(m)
	m.CurState().adjFat(-2, m.MaxState().fat())
	m.Room().ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
}

func (m *Monster) WalkThrough(port *Portal) {
	m.Room().ShowOthers(m, nil, fmt.Sprintf("%s went into a %s.", m.Name(), port.Keyword()))
	m.Room().RemoveMOB(m)
	port.DestRoom().AddMOB(m)
	port.DestRoom().ShowRoom(m)
	m.CurState().adjFat(-2, m.MaxState().fat())
	m.Room().ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
}

func (m *Monster) Room() *Room {
	return m.room
}

func (m *Monster) SetRoom(newRoom *Room) {
	m.room = newRoom
}

func (m *Monster) AttackTarget(target MOB) {
	if target == nil {
		return
	}

	if !m.CurState().alive() {
		return
	} else if !target.CurState().alive() {
		return
	}

	damage := int(m.CurPhyStats().damage()) - int(target.CurPhyStats().defense())
	chance := int(m.CurPhyStats().attack()) - int(target.CurPhyStats().evasion())

	if damage < 0 {
		damage = 0
	}

	m.Room().ShowOthers(m, target, fmt.Sprintf("\r\n%s attacks %s!", m.Name(), target.Name()))

	if chance > 0 {
		target.SendMessage(fmt.Sprintf("%s attacks you!  You are hit for %s damage.",
			m.Name(), CDamageIn(damage)), true)

		target.DamageMOB(m, uint16(damage))
	} else {
		target.SendMessage(fmt.Sprintf("%s attacks you!  They miss!", m.Name()), true)
	}
}

func (m *Monster) KillMOB(killer MOB) {
	m.CurState().setAlive(false)
	m.SetVictim(nil)
	killer.SetVictim(nil)

	expAward := uint64(m.CurPhyStats().level() * 100)
	killer.Room().Show(killer, fmt.Sprintf("%s dies!", m.Name()))
	killer.SendMessage(fmt.Sprintf("You gain %s experience!", CHighlight(expAward)), false)
	killer.AwardExp(expAward)
	m.Room().RemoveMOB(m)

}

func (m *Monster) DamageMOB(attacker MOB, dmg uint16) {
	m.CurState().adjHits(int16(-dmg), m.MaxState().hits())
	if m.CurState().hits() <= 0 {
		m.KillMOB(attacker)
	}

	attacker.SetVictim(m)
	m.SetVictim(attacker)
}

func (m *Monster) Victim() MOB {
	return m.victim
}

func (m *Monster) SetVictim(newVictim MOB) {
	m.victim = newVictim
}

func (m *Monster) AwardExp(howMuch uint64) {
	// TODO(jon) if pet
}

func (m *Monster) AwardRP(howMuch uint32) {
	// TODO(jon) if pet
}
