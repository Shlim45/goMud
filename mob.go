package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MOB struct {
	name          string
	User          *User
	Room          *Room
	CurState      *CharState
	MaxState      *CharState
	BasePhyStats  *PhyStats
	CurPhyStats   *PhyStats
	BaseCharStats *CharStats
	CurCharStats  *CharStats
	Experience    uint64
	RealmPoints   uint32
	inventory     []*Item
	Coins         uint64
	tickType      TickType
	tickCount     uint64
	Victim        *MOB
}

func (m *MOB) isPlayer() bool {
	return m.User != nil
}

func (m *MOB) Inventory() []*Item {
	return m.inventory
}

func (m *MOB) ShowInventory() {
	var inv strings.Builder
	inv.WriteString("You are currently carrying:")
	if len(m.Inventory()) > 0 {
		for _, item := range m.Inventory() {
			inv.WriteString("\r\n  " + CItem(item.FullName()))
		}
	} else {
		inv.WriteString("\r\n  Nothing!")
	}
	m.SendMessage(inv.String(), true)
}

func (m *MOB) AddItem(item *Item) {
	m.inventory = append(m.inventory, item)
	item.SetOwner(m)
}

func (m *MOB) RemoveItem(item *Item) {
	item.SetOwner(nil)

	var items []*Item
	for _, i := range m.inventory {
		if i != item {
			items = append(items, i)
		}
	}
	m.inventory = items
}

func (m *MOB) MoveItemTo(item *Item) {
	if item.Owner() != nil {
		(*item.Owner()).RemoveItem(item)
		item.SetOwner(nil)
	}

	m.inventory = append(m.inventory, item)
	item.SetOwner(m)
}

func (m *MOB) Name() string {
	return m.name
}

func (m *MOB) SetName(newName string) {
	m.name = newName
}

type TickType uint8

const (
	TICK_NORMAL = iota
	TICK_MONSTER
	TICK_PLAYER
	TICK_STOP
)

func (m *MOB) AttackVictim() {
	victim := m.Victim
	if victim == nil {
		return
	}
	if victim.Room == m.Room {
		m.attackTarget(victim)
	}
}

func (m *MOB) Tick(tType TickType) bool {
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
			// hit victim or random opponent
			m.AttackVictim()
		}

	case TICK_PLAYER:
		if m.curState().alive() {
			if m.tickCount%2 == 0 {
				// heal power
				m.adjPower(1, m.maxState().power())
			} else if m.tickCount%7 == 0 {
				// heal hits and fat
				m.adjHits(6, m.maxState().hits())
				m.adjFat(8, m.maxState().fat())
			}
		}

	case TICK_STOP:
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	return m.Tick(m.tickType)
}

func (m *MOB) Init() {
	m.Experience = 0
	rand.Seed(time.Now().UnixMilli())

	baseCStats := CharStats{
		Strength:     uint8(rand.Intn(17) + 4),
		Constitution: uint8(rand.Intn(17) + 4),
		Agility:      uint8(rand.Intn(17) + 4),
		Dexterity:    uint8(rand.Intn(17) + 4),
		Intelligence: uint8(rand.Intn(17) + 4),
		Wisdom:       uint8(rand.Intn(17) + 4),
	}
	m.BaseCharStats = &baseCStats
	m.CurCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(3 * (m.BaseCharStats.dexterity() / 4.0)),
		Damage:       uint16(3 * (m.BaseCharStats.strength() / 4.0)),
		Evasion:      uint16(m.BaseCharStats.agility() / 2.0),
		Defense:      uint16(m.BaseCharStats.constitution() / 2.0),
		MagicAttack:  uint16(m.BaseCharStats.wisdom()),
		MagicDamage:  uint16(m.BaseCharStats.intelligence()),
		MagicEvasion: uint16(3 * (m.BaseCharStats.wisdom() / 4.0)),
		MagicDefense: uint16(3 * (m.BaseCharStats.intelligence() / 4.0)),
		Level:        1 + uint8(m.Experience/1000),
	}
	m.BasePhyStats = &baseStats
	m.CurPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16(30 + (m.BaseCharStats.constitution() / 10.0)),
		Fat:      uint16(30 + (m.BaseCharStats.constitution() / 10.0)),
		Power:    uint16(20 + (m.BaseCharStats.intelligence() / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	m.MaxState = &baseState
	m.CurState = baseState.copyOf()

	if m.tickType == TICK_STOP {
		if m.isPlayer() {
			m.tickType = TICK_PLAYER
		} else {
			m.tickType = TICK_MONSTER
		}
		go m.Tick(m.tickType)
	}
}

func (m *MOB) SendMessage(msg string, newLine bool) {
	if !m.isPlayer() {
		return
	}

	if newLine {
		msg = "\r\n" + msg
	}
	m.User.Session.WriteLine(msg)
}

func (m *MOB) basePhyStats() *PhyStats {
	return m.BasePhyStats
}

func (m *MOB) curPhyStats() *PhyStats {
	return m.CurPhyStats
}

func (m *MOB) recoverPhyStats() {
	m.CurPhyStats = m.BasePhyStats.copyOf()
}

func (m *MOB) curState() *CharState {
	return m.CurState
}

func (m *MOB) maxState() *CharState {
	return m.MaxState
}

func (m *MOB) recoverCharState() {
	m.CurState = m.maxState().copyOf()
}

func (m *MOB) curCharStats() *CharStats {
	return m.CurCharStats
}

func (m *MOB) baseCharStats() *CharStats {
	return m.BaseCharStats
}

func (m *MOB) recoverCharStats() {
	m.CurCharStats = m.BaseCharStats.copyOf()
}

func (m *MOB) adjHits(amount int16, max uint16) {
	newHits := int32(m.curState().hits()) + int32(amount)
	if newHits < 0 {
		m.curState().setHits(0)
	} else if newHits > int32(max) {
		m.curState().setHits(max)
	} else {
		m.curState().setHits(uint16(newHits))
	}
}

func (m *MOB) adjFat(amount int16, max uint16) {
	newFat := int32(m.curState().fat()) + int32(amount)
	if newFat < 0 {
		m.curState().setFat(0)
	} else if newFat > int32(max) {
		m.curState().setFat(max)
	} else {
		m.curState().setFat(uint16(newFat))
	}
}

func (m *MOB) adjPower(amount int16, max uint16) {
	newPower := int32(m.curState().power()) + int32(amount)
	if newPower < 0 {
		m.curState().setPower(0)
	} else if newPower > int32(max) {
		m.curState().setPower(max)
	} else {
		m.curState().setPower(uint16(newPower))
	}
}

func (m *MOB) adjMaxHits(amount int16) {
	newHits := int32(m.maxState().hits()) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	m.maxState().setHits(uint16(newHits))
}

func (m *MOB) adjMaxFat(amount uint16) {
	newFat := int32(m.maxState().fat()) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	m.maxState().setFat(uint16(newFat))
}

func (m *MOB) adjMaxPower(amount uint16) {
	newPower := int32(m.maxState().power()) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	m.maxState().setPower(uint16(newPower))
}

func (m *MOB) attackTarget(target *MOB) {
	if target == nil {
		m.SendMessage("You must specify a target.", true)
	}

	if !m.curState().alive() {
		m.SendMessage("You must be alive to do that!", true)
		return
	} else if !target.curState().alive() {
		m.SendMessage(fmt.Sprintf("%s is already dead!", target.Name()), true)
		return
	}

	damage := int(m.curPhyStats().damage()) - int(target.curPhyStats().defense())
	chance := int(m.curPhyStats().attack()) - int(target.curPhyStats().evasion())

	if damage < 0 {
		damage = 0
	}

	m.Room.ShowOthers(m, target, fmt.Sprintf("\r\n%s attacks %s with their bare hands!", m.Name(), target.Name()))

	if chance > 0 {
		m.SendMessage(fmt.Sprintf("You attack %s with your bare hands and hit for %s damage.",
			target.Name(), CDamageOut(damage)), true)

		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  You are hit for %s damage.",
			m.Name(), CDamageIn(damage)), true)

		target.damageMOB(m, uint16(damage))
	} else {
		m.SendMessage(fmt.Sprintf("You attack %s with your bare hands!  You miss!", target.Name()), true)
		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  They miss!", m.Name()), true)
	}
}

func (m *MOB) AwardExp(howMuch uint64) {
	old := m.Experience
	tnl := 1000 - (old % 1000)
	m.Experience += howMuch
	if howMuch >= tnl {
		newLevel := m.curPhyStats().level() + 1
		if newLevel > 75 {
			return
		}
		m.basePhyStats().setLevel(newLevel)
		m.SendMessage(fmt.Sprintf("You raise a level!\r\n  Your new level is %s.",
			CHighlight(newLevel)), true)
	}
}

func (m *MOB) AwardRP(howMuch uint32) {
	old := m.RealmPoints
	tnr := 100 - (old % 100)
	m.RealmPoints += howMuch
	if howMuch >= tnr {
		m.SendMessage(fmt.Sprintf("You gain a rank in your realm!  Your new title is %s.",
			CHighlight("Dark Acolyte")), true)
	}
}

func (m *MOB) killMOB(killer *MOB) {
	m.curState().setAlive(false)
	m.Victim = nil
	killer.Victim = nil
	if m.isPlayer() {
		m.SendMessage(fmt.Sprintf("You were just killed by %s!", killer.Name()), true)
		killer.SendMessage(fmt.Sprintf("You just killed %s!", m.Name()), true)
		m.Room.ShowOthers(m, killer, fmt.Sprintf("%s was just killed by %s!", m.Name(), killer.Name()))
		// drop held items
		// handle RP
		rpAward := uint32(m.curPhyStats().level() - (killer.curPhyStats().level() - m.curPhyStats().level()))
		plural := "point"
		if rpAward != 1 {
			plural = "points"
		}
		killer.SendMessage(fmt.Sprintf("You gain %s realm %s!", CHighlight(rpAward), plural), false)
		killer.AwardRP(rpAward)
		// create a corpse?  flag
	} else {
		expAward := uint64(m.curPhyStats().level() * 100)
		killer.Room.Show(killer, fmt.Sprintf("%s dies!", m.Name()))
		killer.SendMessage(fmt.Sprintf("You gain %s experience!", CHighlight(expAward)), false)
		killer.AwardExp(expAward)
		m.Room.RemoveMOB(m)
	}
}

func (m *MOB) damageMOB(attacker *MOB, dmg uint16) {
	if m.isPlayer() {
		if (m.curState().hits() == 0) && (dmg > 0) {
			m.killMOB(attacker)
			return
		}
		m.adjHits(int16(-dmg), m.maxState().hits())
		if m.curState().hits() == 0 {
			m.SendMessage("You are almost dead!", true)
		}
	} else {
		m.adjHits(int16(-dmg), m.maxState().hits())
		if m.curState().hits() <= 0 {
			m.killMOB(attacker)
		}
	}

	attacker.Victim = m
	m.Victim = attacker
}

func (m *MOB) recallCorpse(w *World) {
	m.recoverCharState()
	m.recoverPhyStats()
	target := w.GetRoomById("A")
	if target != nil {
		m.SendMessage("You recall your corpse!", true)
		m.Room.ShowOthers(m, nil, fmt.Sprintf("%s recalls their corpse!", m.Name()))
		w.MoveCharacter(m, target)
		m.Room.ShowOthers(m, nil, fmt.Sprintf("%s appears in a puff of smoke.", m.Name()))
		return
	}
}

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}
