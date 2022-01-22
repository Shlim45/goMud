package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Monster struct {
	name          string `json:"name"`
	Room          *Room
	CurState      *CharState
	MaxState      *CharState
	BasePhyStats  *PhyStats
	CurPhyStats   *PhyStats
	BaseCharStats *CharStats
	CurCharStats  *CharStats
	Experience    uint64 `json:"exp"`
	inventory     []*Item
	Coins         uint64 `json:"coins"`
	tickType      TickType
	tickCount     uint64
	Victim        *MOB
}

func (m *Monster) isPlayer() bool {
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

func (m *Monster) AttackVictim() {
	victim := m.Victim
	if victim == nil {
		return
	}
	if victim.Room == m.Room {
		m.attackTarget(victim)
	}
}

func (m *Monster) Tick(tType TickType) bool {
	if m == nil || tType == TICK_STOP {
		return false
	}

	m.tickCount++
	switch m.tickType {
	case TICK_NORMAL:

	case TICK_MONSTER:
		if !m.curState().alive() {
			return false
		}

		if m.tickCount%4 == 0 {
			m.AttackVictim()
		}

	case TICK_STOP:
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	go m.Tick(m.tickType)
	return m.tickType != TICK_STOP
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
	m.BaseCharStats = &baseCStats
	m.CurCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(3 * (m.BaseCharStats.Stats()[STAT_DEXTERITY] / 4)),
		Damage:       uint16(3 * (m.BaseCharStats.Stats()[STAT_STRENGTH] / 4)),
		Evasion:      uint16(m.BaseCharStats.Stats()[STAT_AGILITY] / 2),
		Defense:      uint16(m.BaseCharStats.Stats()[STAT_CONSTITUTION] / 2),
		MagicAttack:  uint16(m.BaseCharStats.Stats()[STAT_WISDOM]),
		MagicDamage:  uint16(m.BaseCharStats.Stats()[STAT_INTELLIGENCE]),
		MagicEvasion: uint16(3 * (m.BaseCharStats.Stats()[STAT_WISDOM] / 4)),
		MagicDefense: uint16(3 * (m.BaseCharStats.Stats()[STAT_INTELLIGENCE] / 4)),
		Level:        1 + uint8(m.Experience/1000),
	}
	m.BasePhyStats = &baseStats
	m.CurPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16((baseStats.Level * 30) + (m.BaseCharStats.Stats()[STAT_CONSTITUTION] / 10.0)),
		Fat:      uint16((baseStats.Level * 30) + (m.BaseCharStats.Stats()[STAT_CONSTITUTION] / 10.0)),
		Power:    uint16((baseStats.Level * 20) + (m.BaseCharStats.Stats()[STAT_INTELLIGENCE] / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	m.MaxState = &baseState
	m.CurState = baseState.copyOf()

	// should only be stop on new mobs
	if m.tickType == TICK_STOP {
		m.Experience = 0
		if m.isPlayer() {
			m.tickType = TICK_PLAYER
		} else {
			m.tickType = TICK_MONSTER
		}
		go m.Tick(m.tickType)
	}
}

func (m *Monster) basePhyStats() *PhyStats {
	return m.BasePhyStats
}

func (m *Monster) curPhyStats() *PhyStats {
	return m.CurPhyStats
}

func (m *Monster) recoverPhyStats() {
	m.CurPhyStats = m.BasePhyStats.copyOf()
}

func (m *Monster) curState() *CharState {
	return m.CurState
}

func (m *Monster) maxState() *CharState {
	return m.MaxState
}

func (m *Monster) recoverCharState() {
	m.CurState = m.maxState().copyOf()
}

func (m *Monster) curCharStats() *CharStats {
	return m.CurCharStats
}

func (m *Monster) baseCharStats() *CharStats {
	return m.BaseCharStats
}

func (m *Monster) recoverCharStats() {
	m.CurCharStats = m.BaseCharStats.copyOf()
}

func (m *Monster) adjHits(amount int16, max uint16) {
	newHits := int32(m.curState().hits()) + int32(amount)
	if newHits < 0 {
		m.curState().setHits(0)
	} else if newHits > int32(max) {
		m.curState().setHits(max)
	} else {
		m.curState().setHits(uint16(newHits))
	}
}

func (m *Monster) adjFat(amount int16, max uint16) {
	newFat := int32(m.curState().fat()) + int32(amount)
	if newFat < 0 {
		m.curState().setFat(0)
	} else if newFat > int32(max) {
		m.curState().setFat(max)
	} else {
		m.curState().setFat(uint16(newFat))
	}
}

func (m *Monster) adjPower(amount int16, max uint16) {
	newPower := int32(m.curState().power()) + int32(amount)
	if newPower < 0 {
		m.curState().setPower(0)
	} else if newPower > int32(max) {
		m.curState().setPower(max)
	} else {
		m.curState().setPower(uint16(newPower))
	}
}

func (m *Monster) adjMaxHits(amount int16) {
	newHits := int32(m.maxState().hits()) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	m.maxState().setHits(uint16(newHits))
}

func (m *Monster) adjMaxFat(amount uint16) {
	newFat := int32(m.maxState().fat()) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	m.maxState().setFat(uint16(newFat))
}

func (m *Monster) adjMaxPower(amount uint16) {
	newPower := int32(m.maxState().power()) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	m.maxState().setPower(uint16(newPower))
}

func (m *Monster) Walk(dest *Room, verb string) {
	m.Room.ShowOthers(m, nil, fmt.Sprintf("%s went %s.", m.Name(), verb))
	m.Room.RemoveMOB(m)
	dest.AddMOB(m)
	dest.ShowRoom(m)
	m.adjFat(-2, m.maxState().fat())
	m.Room.ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
}

func (m *Monster) WalkThrough(port *Portal) {
	m.Room.ShowOthers(m, nil, fmt.Sprintf("%s went into a %s.", m.Name(), port.Keyword()))
	m.Room.RemoveMOB(m)
	port.DestRoom().AddMOB(m)
	port.DestRoom().ShowRoom(m)
	m.adjFat(-2, m.maxState().fat())
	m.Room.ShowOthers(m, nil, fmt.Sprintf("%s just came in.", m.Name()))
}

func (m *Monster) attackTarget(target *MOB) {
	if target == nil {
		return
	}

	if !m.curState().alive() {
		return
	} else if !target.curState().alive() {
		return
	}

	damage := int(m.curPhyStats().damage()) - int(target.curPhyStats().defense())
	chance := int(m.curPhyStats().attack()) - int(target.curPhyStats().evasion())

	if damage < 0 {
		damage = 0
	}

	m.Room.ShowOthers(m, target, fmt.Sprintf("\r\n%s attacks %s!", m.Name(), target.Name()))

	if chance > 0 {
		target.SendMessage(fmt.Sprintf("%s attacks you!  You are hit for %s damage.",
			m.Name(), CDamageIn(damage)), true)

		target.damageMOB(m, uint16(damage))
	} else {
		target.SendMessage(fmt.Sprintf("%s attacks you!  They miss!", m.Name()), true)
	}
}

func (m *Monster) killMOB(killer *MOB) {
	m.curState().setAlive(false)
	m.Victim = nil
	killer.Victim = nil

	expAward := uint64(m.curPhyStats().level() * 100)
	killer.Room.Show(killer, fmt.Sprintf("%s dies!", m.Name()))
	killer.SendMessage(fmt.Sprintf("You gain %s experience!", CHighlight(expAward)), false)
	killer.AwardExp(expAward)
	m.Room.RemoveMOB(m)

}

func (m *Monster) damageMOB(attacker *MOB, dmg uint16) {
	m.adjHits(int16(-dmg), m.maxState().hits())
	if m.curState().hits() <= 0 {
		m.killMOB(attacker)
	}

	attacker.Victim = m
	m.Victim = attacker
}
