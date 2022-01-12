package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Player struct {
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
	inventory     []*Item
	Coins         uint64
}

func (p *Player) Inventory() []*Item {
	return p.inventory
}

func (p *Player) ShowInventory() {
	var inv strings.Builder
	inv.WriteString("You are currently carrying:")
	if len(p.Inventory()) > 0 {
		for _, item := range p.Inventory() {
			inv.WriteString("\r\n  " + CItem(item.FullName()))
		}
	} else {
		inv.WriteString("\r\n  Nothing!")
	}
	p.SendMessage(inv.String(), true)
}

func (p *Player) AddItem(item *Item) {
	p.inventory = append(p.inventory, item)
	item.SetOwner(p)
}

func (p *Player) RemoveItem(item *Item) {
	item.SetOwner(nil)

	var items []*Item
	for _, i := range p.inventory {
		if i != item {
			items = append(items, i)
		}
	}
	p.inventory = items
}

func (p *Player) MoveItemTo(item *Item) {
	if item.Owner() != nil {
		(*item.Owner()).RemoveItem(item)
		item.SetOwner(nil)
	}

	p.inventory = append(p.inventory, item)
	item.SetOwner(p)
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) SetName(newName string) {
	p.name = newName
}

func (p *Player) Init() {
	p.Experience = 0
	rand.Seed(time.Now().UnixMilli())

	baseCStats := CharStats{
		Strength:     uint8(rand.Intn(17) + 4),
		Constitution: uint8(rand.Intn(17) + 4),
		Agility:      uint8(rand.Intn(17) + 4),
		Dexterity:    uint8(rand.Intn(17) + 4),
		Intelligence: uint8(rand.Intn(17) + 4),
		Wisdom:       uint8(rand.Intn(17) + 4),
	}
	p.BaseCharStats = &baseCStats
	p.CurCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(5 * (p.BaseCharStats.dexterity() / 10.0)),
		Damage:       uint16(5 * (p.BaseCharStats.strength() / 10.0)),
		Evasion:      uint16(5 * (p.BaseCharStats.agility() / 10.0)),
		Defense:      uint16(5 * (p.BaseCharStats.constitution() / 10.0)),
		MagicAttack:  uint16(5 * (p.BaseCharStats.wisdom() / 10.0)),
		MagicDamage:  uint16(5 * (p.BaseCharStats.intelligence() / 10.0)),
		MagicEvasion: uint16(3 * (p.BaseCharStats.wisdom() / 4.0)),
		MagicDefense: uint16(3 * (p.BaseCharStats.intelligence() / 4.0)),
		Level:        1 + uint8(p.Experience/1000),
	}
	p.BasePhyStats = &baseStats
	p.CurPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16(30 * (p.BaseCharStats.constitution() / 10.0)),
		Fat:      uint16(30 * (p.BaseCharStats.constitution() / 10.0)),
		Power:    uint16(20 * (p.BaseCharStats.intelligence() / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	p.MaxState = &baseState
	p.CurState = baseState.copyOf()
}

func (p *Player) SendMessage(msg string, newLine bool) {
	if newLine {
		msg = "\r\n" + msg
	}
	p.User.Session.WriteLine(msg)
}

func (p *Player) basePhyStats() *PhyStats {
	return p.BasePhyStats
}

func (p *Player) curPhyStats() *PhyStats {
	return p.CurPhyStats
}

func (p *Player) recoverPhyStats() {
	p.CurPhyStats = p.BasePhyStats.copyOf()
}

func (p *Player) curState() *CharState {
	return p.CurState
}

func (p *Player) maxState() *CharState {
	return p.MaxState
}

func (p *Player) recoverCharState() {
	p.CurState = p.maxState().copyOf()
}

func (p *Player) curCharStats() *CharStats {
	return p.CurCharStats
}

func (p *Player) baseCharStats() *CharStats {
	return p.BaseCharStats
}

func (p *Player) recoverCharStats() {
	p.CurCharStats = p.BaseCharStats.copyOf()
}

func (p *Player) level() uint8 {
	return p.curPhyStats().level()
}

func (p *Player) setLevel(newLevel uint8) {
	p.basePhyStats().setLevel(newLevel)
	p.recoverPhyStats()
}

func (p *Player) adjHits(amount int16) {
	newHits := int32(p.curState().hits()) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	p.curState().setHits(uint16(newHits))
}

func (p *Player) adjFat(amount int16) {
	newFat := int32(p.curState().fat()) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	p.curState().setFat(uint16(newFat))
}

func (p *Player) adjPower(amount int16) {
	newPower := int32(p.curState().power()) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	p.curState().setPower(uint16(newPower))
}

func (p *Player) adjMaxHits(amount int16) {
	newHits := int32(p.maxState().hits()) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	p.maxState().setHits(uint16(newHits))
}

func (p *Player) adjMaxFat(amount uint16) {
	newFat := int32(p.maxState().fat()) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	p.maxState().setFat(uint16(newFat))
}

func (p *Player) adjMaxPower(amount uint16) {
	newPower := int32(p.maxState().power()) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	p.maxState().setPower(uint16(newPower))
}

func (p *Player) attackTarget(target *Player) {
	if target == nil {
		p.SendMessage("You must specify a target.", true)
	}
	if !p.curState().alive() {
		p.SendMessage("You must be alive to do that!", true)
		return
	} else if !target.curState().alive() {
		p.SendMessage(fmt.Sprintf("%s is already dead!", target.Name()), true)
		return
	}

	damage := int(p.curPhyStats().damage()) - int(target.curPhyStats().defense())
	chance := int(p.curPhyStats().attack()) - int(target.curPhyStats().evasion())

	if damage < 0 {
		damage = 0
	}

	if chance > 0 {
		outDamage := CDamageOut(fmt.Sprintf("%d", damage))
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands and hit for %s damage.",
			target.Name(), outDamage), true)

		inDamage := CDamageIn(fmt.Sprintf("%d", damage))
		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  You are hit for %s damage.",
			p.Name(), inDamage), true)

		p.Room.ShowOthers(p, target, fmt.Sprintf("%s attacks %s with their bare hands!\r\n", p.Name(), target.Name()))
		target.damagePlayer(uint16(damage))
	} else {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands!  You miss!", target.Name()), true)
	}
}

func (p *Player) damagePlayer(dmg uint16) {
	if (p.curState().hits() == 0) && (dmg > 0) {
		p.curState().setAlive(false)
		p.SendMessage("You were just killed!", true)
		p.Room.ShowOthers(p, nil, fmt.Sprintf("%s was just killed!", p.Name()))
		return
	}

	p.adjHits(int16(-dmg))
	if p.curState().hits() == 0 {
		p.SendMessage("You are almost dead!", true)
	}
}

func (p *Player) recallCorpse(w *World) {
	p.recoverCharState()
	p.recoverPhyStats()
	target := w.GetRoomById("A")
	if target != nil {
		p.SendMessage("You recall your corpse!", true)
		p.Room.ShowOthers(p, nil, fmt.Sprintf("%s recalls their corpse!", p.Name()))
		w.MoveCharacter(p, target)
		p.Room.ShowOthers(p, nil, fmt.Sprintf("%s appears in a puff of smoke.", p.Name()))
		return
	}
}

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}
