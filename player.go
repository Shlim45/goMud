package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Player struct {
	name          string `json:"name"`
	Account       string `json:"account"`
	User          *User
	Room          *Room
	CurState      *CharState
	MaxState      *CharState
	BasePhyStats  *PhyStats
	CurPhyStats   *PhyStats
	BaseCharStats *CharStats
	CurCharStats  *CharStats
	Experience    uint64 `json:"exp"`
	RealmPoints   uint32 `json:"rp"`
	inventory     []*Item
	Coins         uint64 `json:"coins"`
	tickType      TickType
	tickCount     uint64
	Victim        *MOB
	LastDate      string
}

func (p *Player) isPlayer() bool {
	return true
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
		item.Owner().RemoveItem(item)
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

func (p *Player) AttackVictim() {
	victim := p.Victim
	if victim == nil {
		return
	}
	if victim.Room == p.Room {
		p.attackTarget(victim)
	}
}

func (p *Player) Tick(tType TickType) bool {
	if p == nil || tType == TICK_STOP {
		return false
	}

	p.tickCount++
	switch p.tickType {
	case TICK_NORMAL:

	case TICK_PLAYER:
		if p.curState().alive() {
			if p.tickCount%2 == 0 {
				// heal power
				p.adjPower(1, p.maxState().power())
			} else if p.tickCount%7 == 0 {
				// heal hits and fat
				p.adjHits(6, p.maxState().hits())
				p.adjFat(8, p.maxState().fat())
			}
		}

	case TICK_STOP:
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	go p.Tick(p.tickType)
	return p.tickType != TICK_STOP
}

func (p *Player) Init(library *MudLib) {
	rand.Seed(time.Now().UnixMilli())

	var stats [NUM_STATS]uint8

	var newRace Race
	var newClass CharClass

	for i := range stats {
		stats[i] = uint8(rand.Intn(17) + 4)
	}
	newRace = library.FindRace("Human")
	newClass = library.FindCharClass("Fighter")

	baseCStats := CharStats{
		stats:     stats,
		charClass: newClass,
		race:      newRace,
	}
	p.BaseCharStats = &baseCStats
	p.CurCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(3 * (p.BaseCharStats.Stats()[STAT_DEXTERITY] / 4)),
		Damage:       uint16(3 * (p.BaseCharStats.Stats()[STAT_STRENGTH] / 4)),
		Evasion:      uint16(p.BaseCharStats.Stats()[STAT_AGILITY] / 2),
		Defense:      uint16(p.BaseCharStats.Stats()[STAT_CONSTITUTION] / 2),
		MagicAttack:  uint16(p.BaseCharStats.Stats()[STAT_WISDOM]),
		MagicDamage:  uint16(p.BaseCharStats.Stats()[STAT_INTELLIGENCE]),
		MagicEvasion: uint16(3 * (p.BaseCharStats.Stats()[STAT_WISDOM] / 4)),
		MagicDefense: uint16(3 * (p.BaseCharStats.Stats()[STAT_INTELLIGENCE] / 4)),
		Level:        1 + uint8(p.Experience/1000),
	}
	p.BasePhyStats = &baseStats
	p.CurPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16((baseStats.Level * 30) + (p.BaseCharStats.Stats()[STAT_CONSTITUTION] / 10.0)),
		Fat:      uint16((baseStats.Level * 30) + (p.BaseCharStats.Stats()[STAT_CONSTITUTION] / 10.0)),
		Power:    uint16((baseStats.Level * 20) + (p.BaseCharStats.Stats()[STAT_INTELLIGENCE] / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	p.MaxState = &baseState
	p.CurState = baseState.copyOf()

	// should only be stop on new mobs
	if p.tickType == TICK_STOP {
		p.Experience = 0
		if p.isPlayer() {
			p.tickType = TICK_PLAYER
		} else {
			p.tickType = TICK_MONSTER
		}
		go p.Tick(p.tickType)
	}
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

func (p *Player) adjHits(amount int16, max uint16) {
	newHits := int32(p.curState().hits()) + int32(amount)
	if newHits < 0 {
		p.curState().setHits(0)
	} else if newHits > int32(max) {
		p.curState().setHits(max)
	} else {
		p.curState().setHits(uint16(newHits))
	}
}

func (p *Player) adjFat(amount int16, max uint16) {
	newFat := int32(p.curState().fat()) + int32(amount)
	if newFat < 0 {
		p.curState().setFat(0)
	} else if newFat > int32(max) {
		p.curState().setFat(max)
	} else {
		p.curState().setFat(uint16(newFat))
	}
}

func (p *Player) adjPower(amount int16, max uint16) {
	newPower := int32(p.curState().power()) + int32(amount)
	if newPower < 0 {
		p.curState().setPower(0)
	} else if newPower > int32(max) {
		p.curState().setPower(max)
	} else {
		p.curState().setPower(uint16(newPower))
	}
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

func (p *Player) Walk(dest *Room, verb string) {
	p.SendMessage(fmt.Sprintf("You travel %s.", verb), true)
	p.Room.ShowOthers(p, nil, fmt.Sprintf("%s went %s.", p.Name(), verb))
	p.Room.RemoveMOB(p)
	dest.AddMOB(p)
	dest.ShowRoom(p)
	p.adjFat(-2, p.maxState().fat())
	p.Room.ShowOthers(p, nil, fmt.Sprintf("%s just came in.", p.Name()))
}

func (p *Player) WalkThrough(port *Portal) {
	p.SendMessage(fmt.Sprintf("You travel into a %s.", port.Keyword()), true)
	p.Room.ShowOthers(p, nil, fmt.Sprintf("%s went into a %s.", p.Name(), port.Keyword()))
	p.Room.RemoveMOB(p)
	port.DestRoom().AddMOB(p)
	port.DestRoom().ShowRoom(p)
	p.adjFat(-2, p.maxState().fat())
	p.Room.ShowOthers(p, nil, fmt.Sprintf("%s just came in.", p.Name()))
}

func (p *Player) attackTarget(target *MOB) {
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

	p.Room.ShowOthers(p, target, fmt.Sprintf("\r\n%s attacks %s with their bare hands!", p.Name(), target.Name()))

	if chance > 0 {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands and hit for %s damage.",
			target.Name(), CDamageOut(damage)), true)

		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  You are hit for %s damage.",
			p.Name(), CDamageIn(damage)), true)

		target.damageMOB(p, uint16(damage))
	} else {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands!  You miss!", target.Name()), true)
		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  They miss!", p.Name()), true)
	}
}

func (p *Player) AwardExp(howMuch uint64) {
	old := p.Experience
	tnl := 1000 - (old % 1000)
	p.Experience += howMuch
	if howMuch >= tnl {
		newLevel := 1 + uint8(p.Experience/1000)
		if newLevel > 75 {
			return
		}
		p.basePhyStats().setLevel(newLevel)
		p.recoverPhyStats()
		p.SendMessage(fmt.Sprintf("You raise a level!\r\n  Your new level is %s.",
			CHighlight(newLevel)), true)
	}
}

func (p *Player) AwardRP(howMuch uint32) {
	old := p.RealmPoints
	tnr := 100 - (old % 100)
	p.RealmPoints += howMuch
	if howMuch >= tnr {
		p.SendMessage(fmt.Sprintf("You gain a rank in your realm!  Your new title is %s.",
			CHighlight("Dark Acolyte")), true)
	}
}

func (p *Player) killMOB(killer *MOB) {
	p.curState().setAlive(false)
	p.Victim = nil
	killer.Victim = nil
	if p.isPlayer() {
		p.SendMessage(fmt.Sprintf("You were just killed by %s!", killer.Name()), true)
		killer.SendMessage(fmt.Sprintf("You just killed %s!", p.Name()), true)
		p.Room.ShowOthers(p, killer, fmt.Sprintf("%s was just killed by %s!", p.Name(), killer.Name()))
		// drop held items
		// handle RP
		rpAward := uint32(p.curPhyStats().level() - (killer.curPhyStats().level() - p.curPhyStats().level()))
		plural := "point"
		if rpAward != 1 {
			plural = "points"
		}
		killer.SendMessage(fmt.Sprintf("You gain %s realm %s!", CHighlight(rpAward), plural), false)
		killer.AwardRP(rpAward)
		// create a corpse?  flag
	} else {
		expAward := uint64(p.curPhyStats().level() * 100)
		killer.Room.Show(killer, fmt.Sprintf("%s dies!", p.Name()))
		killer.SendMessage(fmt.Sprintf("You gain %s experience!", CHighlight(expAward)), false)
		killer.AwardExp(expAward)
		p.Room.RemoveMOB(p)
	}
}

func (p *Player) damageMOB(attacker *MOB, dmg uint16) {
	if p.isPlayer() {
		if (p.curState().hits() == 0) && (dmg > 0) {
			p.killMOB(attacker)
			return
		}
		p.adjHits(int16(-dmg), p.maxState().hits())
		if p.curState().hits() == 0 {
			p.SendMessage("You are almost dead!", true)
		}
	} else {
		p.adjHits(int16(-dmg), p.maxState().hits())
		if p.curState().hits() <= 0 {
			p.killMOB(attacker)
		}
	}

	attacker.Victim = p
	p.Victim = attacker
}

// TODO(jon): Player only
func (p *Player) releaseCorpse(w *World) {
	p.recoverCharState()
	p.recoverPhyStats()
	target := w.GetRoomById("A")
	if target != nil {
		p.SendMessage("You release your corpse!", true)
		p.Room.ShowOthers(p, nil, fmt.Sprintf("%s releases their corpse!", p.Name()))
		w.MoveCharacter(p, target)
		p.Room.ShowOthers(p, nil, fmt.Sprintf("%s appears in a puff of smoke.", p.Name()))
		return
	}
}

type PlayerDB struct {
	name       string `json:"name"`
	account    string `json:"account"`
	class      string `json:"class"`
	race       string `json:"race"`
	room       string `json:"room"`
	coins      uint64 `json:"coins"`
	stre       uint8  `json:"stre"`
	cons       uint8  `json:"cons"`
	agil       uint8  `json:"agil"`
	dext       uint8  `json:"dext"`
	inte       uint8  `json:"inte"`
	wisd       uint8  `json:"wisd"`
	con_loss   uint8  `json:"con_loss"`
	level      uint8  `json:"level"`
	exp        uint64 `json:"exp"`
	rp         uint32 `json:"rp"`
	hits       uint16 `json:"hits"`
	fat        uint16 `json:"fat"`
	power      uint16 `json:"power"`
	trains     uint16 `json:"trains"`
	guild      string `json:"guild"`
	guild_rank uint8  `json:"guild_rank"`
	last_date  string `json:"last_date"`
}

func (p *Player) SavePlayerToDBQuery() (string, error) {
	if p == nil {
		return "", errors.New("MOB is nil")
	}
	var location string
	if p.Room != nil {
		location = p.Room.RoomID()
	} else {
		// TODO(jon): need to preserve last room
		location = ""
	}
	player := PlayerDB{
		name:       p.Name(),
		account:    p.Account,
		class:      p.baseCharStats().CurrentClass().Name(),
		race:       p.baseCharStats().Race().Name(),
		room:       location,
		coins:      p.Coins,
		stre:       p.baseCharStats().Stat(STAT_STRENGTH),
		cons:       p.baseCharStats().Stat(STAT_CONSTITUTION),
		agil:       p.baseCharStats().Stat(STAT_AGILITY),
		dext:       p.baseCharStats().Stat(STAT_DEXTERITY),
		inte:       p.baseCharStats().Stat(STAT_INTELLIGENCE),
		wisd:       p.baseCharStats().Stat(STAT_WISDOM),
		con_loss:   0,
		level:      p.basePhyStats().level(),
		exp:        p.Experience,
		rp:         p.RealmPoints,
		hits:       p.maxState().hits(),
		fat:        p.maxState().fat(),
		power:      p.maxState().power(),
		trains:     0,
		guild:      "",
		guild_rank: 0,
		last_date:  p.LastDate,
	}
	return fmt.Sprintf("INSERT INTO Player VALUES ('%s', '%s', '%s', '%s', '%s', %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, '%s', %d, '%s')"+
			" AS new ON DUPLICATE KEY UPDATE name=new.name, account=new.account, class=new.class, race=new.race, room=new.room, coins=new.coins, stre=new.stre, cons=new.cons, agil=new.agil, dext=new.dext, inte=new.inte, wisd=new.wisd, "+
			"con_loss=new.con_loss, level=new.level, exp=new.exp, rp=new.rp, hits=new.hits, fat=new.fat, power=new.power, trains=new.trains, guild=new.guild, guild_rank=new.guild_rank, last_date=new.last_date",
			player.name, player.account, player.class, player.race, player.room, player.coins, player.stre, player.cons, player.agil, player.dext, player.inte, player.wisd,
			player.con_loss, player.level, player.exp, player.rp, player.hits, player.fat, player.power, player.trains, player.guild, player.guild_rank, player.last_date),
		nil // the error
}

func CreatePlayerTableDBQuery() string {
	return "CREATE TABLE IF NOT EXISTS Player(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"account VARCHAR(20)," +
		"class VARCHAR(20)," +
		"race VARCHAR(20)," +
		"room VARCHAR(60)," +
		"coins BIGINT UNSIGNED NOT NULL," +
		"stre TINYINT UNSIGNED NOT NULL," +
		"cons TINYINT UNSIGNED NOT NULL," +
		"agil TINYINT UNSIGNED NOT NULL," +
		"dext TINYINT UNSIGNED NOT NULL," +
		"inte TINYINT UNSIGNED NOT NULL," +
		"wisd TINYINT UNSIGNED NOT NULL," +
		"con_loss TINYINT UNSIGNED NOT NULL," +
		"level TINYINT UNSIGNED NOT NULL," +
		"exp BIGINT UNSIGNED NOT NULL," +
		"rp INT UNSIGNED NOT NULL," +
		"hits SMALLINT UNSIGNED NOT NULL," +
		"fat SMALLINT UNSIGNED NOT NULL," +
		"power SMALLINT UNSIGNED NOT NULL," +
		"trains SMALLINT UNSIGNED NOT NULL DEFAULT 0," +
		"guild VARCHAR(30)," + // FK
		"guild_rank TINYINT UNSIGNED NOT NULL DEFAULT 0," +
		"last_date TIMESTAMP," +
		"FOREIGN KEY (account) REFERENCES Account(username) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (class) REFERENCES CharClass(name) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (race) REFERENCES Race(name) ON UPDATE CASCADE ON DELETE SET NULL" +
		")"
}
