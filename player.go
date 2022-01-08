package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	Name         string
	User         *User
	Room         *Room
	CurState     *CharState
	MaxState     *CharState
	BasePhyStats *PhyStats
	CurPhyStats  *PhyStats
	Experience   uint32
}

func (p *Player) Init() {
	p.Experience = 0
	rand.Seed(time.Now().UnixMilli())
	baseStats := PhyStats{
		Attack:       uint16(rand.Intn(11)),
		Damage:       uint16(rand.Intn(11)),
		Evasion:      uint16(rand.Intn(11)),
		Defense:      uint16(rand.Intn(11)),
		MagicAttack:  uint16(rand.Intn(11)),
		MagicDamage:  uint16(rand.Intn(11)),
		MagicEvasion: uint16(rand.Intn(11)),
		MagicDefense: uint16(rand.Intn(11)),
		Level:        1 + uint8(p.Experience/1000),
	}
	p.BasePhyStats = &baseStats
	p.CurPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     30,
		Fat:      30,
		Power:    20,
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
	p.CurState = p.MaxState.copyOf()
}

func (p *Player) level() uint8 {
	return p.curPhyStats().level()
}

func (p *Player) setLevel(newLevel uint8) {
	p.basePhyStats().setLevel(newLevel)
	p.recoverPhyStats()
}

func (p *Player) adjHits(amount int16) {
	newHits := int32(p.CurState.Hits) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	p.CurState.Hits = uint16(newHits)
}

func (p *Player) adjFat(amount int16) {
	newFat := int32(p.CurState.Fat) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	p.CurState.Fat = uint16(newFat)
}

func (p *Player) adjPower(amount int16) {
	newPower := int32(p.CurState.Power) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	p.CurState.Power = uint16(newPower)
}

func (p *Player) adjMaxHits(amount int16) {
	newHits := int32(p.MaxState.Hits) + int32(amount)
	if newHits < 0 {
		newHits = 0
	}
	p.MaxState.Hits = uint16(newHits)
}

func (p *Player) adjMaxFat(amount uint16) {
	newFat := int32(p.MaxState.Fat) + int32(amount)
	if newFat < 0 {
		newFat = 0
	}
	p.MaxState.Fat = uint16(newFat)
}

func (p *Player) adjMaxPower(amount uint16) {
	newPower := int32(p.MaxState.Power) + int32(amount)
	if newPower < 0 {
		newPower = 0
	}
	p.MaxState.Power = uint16(newPower)
}

func (p *Player) attackTarget(target *Player) {
	if target == nil {
		p.SendMessage("You must specify a target.", true)
	}
	if !p.curState().Alive {
		p.SendMessage("You must be alive to do that!", true)
		return
	} else if !target.curState().Alive {
		p.SendMessage(fmt.Sprintf("%s is already dead!", target.Name), true)
		return
	}

	damage := int(p.curPhyStats().Damage) - int(target.curPhyStats().Defense)
	chance := int(p.curPhyStats().Attack) - int(target.curPhyStats().Evasion)

	if damage < 0 {
		damage = 0
	}

	if chance > 0 {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands and hit for %d damage.",
			target.Name, damage), true)
		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  You are hit for %d damage.",
			p.Name, damage), true)
		p.Room.ShowOthers(p, fmt.Sprintf("%s attacks %s with their bare hands!\r\n", p.Name, target.Name))
		target.damagePlayer(uint16(damage))
	} else {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands!  You miss!", target.Name), true)
	}
}

func (p *Player) damagePlayer(dmg uint16) {
	if (p.curState().Hits == 0) && (dmg > 0) {
		p.CurState.Alive = false
		p.SendMessage("You were just killed!", true)
		p.Room.ShowOthers(p, fmt.Sprintf("%s was just killed!", p.Name))
		return
	}

	p.adjHits(int16(-dmg))
	if p.curState().Hits == 0 {
		p.SendMessage("You are almost dead!", true)
	}
}

func (p *Player) recallCorpse(w *World) {
	p.recoverCharState()
	p.recoverPhyStats()
	target := w.GetRoomById("A")
	if target != nil {
		p.SendMessage("You recall your corpse!", true)
		p.Room.ShowOthers(p, fmt.Sprintf("%s recalls their corpse!", p.Name))
		w.MoveCharacter(p, target)
		p.Room.ShowOthers(p, fmt.Sprintf("%s appears in a puff of smoke.", p.Name))
		return
	}
}
