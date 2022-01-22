package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Player struct {
	name          string
	Account       string
	User          *User
	room          *Room
	curState      *CharState
	maxState      *CharState
	basePhyStats  *PhyStats
	curPhyStats   *PhyStats
	baseCharStats *CharStats
	curCharStats  *CharStats
	Experience    uint64
	RealmPoints   uint32
	inventory     []*Item
	coins         uint64
	tickType      TickType
	tickCount     uint64
	victim        MOB
	LastDate      string
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) SetName(newName string) {
	p.name = newName
}

func (p *Player) Coins() uint64 {
	return p.coins
}

func (p *Player) SetCoins(amount uint64) {
	p.coins = amount
}

func (p *Player) AdjCoins(amount int64) {
	if amount < 0 {
		less := uint64(-amount)
		if less > p.coins {
			p.coins = 0
		} else {
			p.coins -= less
		}
	} else {
		p.coins += uint64(amount)
	}
}

func (p *Player) BasePhyStats() *PhyStats {
	return p.basePhyStats
}

func (p *Player) CurPhyStats() *PhyStats {
	return p.curPhyStats
}

func (p *Player) RecoverPhyStats() {
	p.curPhyStats = p.BasePhyStats().copyOf()
}

func (p *Player) CurState() *CharState {
	return p.curState
}

func (p *Player) MaxState() *CharState {
	return p.maxState
}

func (p *Player) RecoverCharState() {
	p.curState = p.MaxState().copyOf()
}

func (p *Player) CurCharStats() *CharStats {
	return p.curCharStats
}

func (p *Player) BaseCharStats() *CharStats {
	return p.baseCharStats
}

func (p *Player) RecoverCharStats() {
	p.curCharStats = p.BaseCharStats().copyOf()
}

func (p *Player) IsPlayer() bool {
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

func (p *Player) AttackVictim() {
	victim := p.Victim()
	if victim == nil {
		return
	}
	if victim.Room() == p.Room() {
		p.AttackTarget(victim)
	}
}

func (p *Player) Tick(tType TickType) bool {
	if p == nil || tType == TickStop {
		return false
	}
	p.tickType = tType
	p.tickCount++
	switch p.tickType {
	case TickNormal:

	case TickPlayer:
		if p.User != nil && p.User.Session.Status() == INGAME && p.CurState().alive() {
			if p.tickCount%2 == 0 {
				// heal power
				p.CurState().adjPower(1, p.MaxState().power())
			} else if p.tickCount%7 == 0 {
				// heal hits and fat
				p.CurState().adjHits(6, p.MaxState().hits())
				p.CurState().adjFat(8, p.MaxState().fat())
			}
		}

	case TickStop:
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	go p.Tick(p.tickType)
	return p.tickType != TickStop
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
	p.baseCharStats = &baseCStats
	p.curCharStats = baseCStats.copyOf()

	baseStats := PhyStats{
		Attack:       uint16(3 * (p.BaseCharStats().Stats()[STAT_DEXTERITY] / 4)),
		Damage:       uint16(3 * (p.BaseCharStats().Stats()[STAT_STRENGTH] / 4)),
		Evasion:      uint16(p.BaseCharStats().Stats()[STAT_AGILITY] / 2),
		Defense:      uint16(p.BaseCharStats().Stats()[STAT_CONSTITUTION] / 2),
		MagicAttack:  uint16(p.BaseCharStats().Stats()[STAT_WISDOM]),
		MagicDamage:  uint16(p.BaseCharStats().Stats()[STAT_INTELLIGENCE]),
		MagicEvasion: uint16(3 * (p.BaseCharStats().Stats()[STAT_WISDOM] / 4)),
		MagicDefense: uint16(3 * (p.BaseCharStats().Stats()[STAT_INTELLIGENCE] / 4)),
		Level:        1 + uint8(p.Experience/1000),
	}
	p.basePhyStats = &baseStats
	p.curPhyStats = baseStats.copyOf()

	baseState := CharState{
		Hits:     uint16((baseStats.Level * 30) + (p.BaseCharStats().Stats()[STAT_CONSTITUTION] / 10.0)),
		Fat:      uint16((baseStats.Level * 30) + (p.BaseCharStats().Stats()[STAT_CONSTITUTION] / 10.0)),
		Power:    uint16((baseStats.Level * 20) + (p.BaseCharStats().Stats()[STAT_INTELLIGENCE] / 10.0)),
		Alive:    true,
		Standing: true,
		Sitting:  false,
		Laying:   false,
	}
	p.maxState = &baseState
	p.curState = baseState.copyOf()

	// should only be stop on new players
	if p.tickType == TickStop {
		p.Experience = 0
		p.tickType = TickPlayer

		go p.Tick(p.tickType)
	}
}

func (p *Player) SendMessage(msg string, newLine bool) {
	if newLine {
		msg = "\r\n" + msg
	}
	p.User.Session.WriteLine(msg)
}

func (p *Player) Walk(dest *Room, verb string) {
	p.SendMessage(fmt.Sprintf("You travel %s.", verb), true)
	p.Room().ShowOthers(p, nil, fmt.Sprintf("%s went %s.", p.Name(), verb))
	p.Room().RemoveMOB(p)
	dest.AddMOB(p)
	dest.ShowRoom(p)
	p.CurState().adjFat(-2, p.MaxState().fat())
	p.Room().ShowOthers(p, nil, fmt.Sprintf("%s just came in.", p.Name()))
}

func (p *Player) WalkThrough(port *Portal) {
	p.SendMessage(fmt.Sprintf("You travel into a %s.", port.Keyword()), true)
	p.Room().ShowOthers(p, nil, fmt.Sprintf("%s went into a %s.", p.Name(), port.Keyword()))
	p.Room().RemoveMOB(p)
	port.DestRoom().AddMOB(p)
	port.DestRoom().ShowRoom(p)
	p.CurState().adjFat(-2, p.MaxState().fat())
	p.Room().ShowOthers(p, nil, fmt.Sprintf("%s just came in.", p.Name()))
}

func (p *Player) Room() *Room {
	return p.room
}

func (p *Player) SetRoom(newRoom *Room) {
	p.room = newRoom
}

func (p *Player) AttackTarget(target MOB) {
	if target == nil {
		p.SendMessage("You must specify a target.", true)
	}

	if !p.CurState().alive() {
		p.SendMessage("You must be alive to do that!", true)
		return
	} else if !target.CurState().alive() {
		p.SendMessage(fmt.Sprintf("%s is already dead!", target.Name()), true)
		return
	}

	damage := int(p.CurPhyStats().damage()) - int(target.CurPhyStats().defense())
	chance := int(p.CurPhyStats().attack()) - int(target.CurPhyStats().evasion())

	if damage < 0 {
		damage = 0
	}

	p.Room().ShowOthers(p, target, fmt.Sprintf("\r\n%s attacks %s with their bare hands!", p.Name(), target.Name()))

	if chance > 0 {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands and hit for %s damage.",
			target.Name(), CDamageOut(damage)), true)

		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  You are hit for %s damage.",
			p.Name(), CDamageIn(damage)), true)

		target.DamageMOB(p, uint16(damage))
	} else {
		p.SendMessage(fmt.Sprintf("You attack %s with your bare hands!  You miss!", target.Name()), true)
		target.SendMessage(fmt.Sprintf("%s attacks you with their bare hands!  They miss!", p.Name()), true)
	}
}

func (p *Player) KillMOB(killer MOB) {
	p.CurState().setAlive(false)
	p.SetVictim(nil)
	killer.SetVictim(nil)

	p.SendMessage(fmt.Sprintf("You were just killed by %s!", killer.Name()), true)
	killer.SendMessage(fmt.Sprintf("You just killed %s!", p.Name()), true)
	p.Room().ShowOthers(p, killer, fmt.Sprintf("%s was just killed by %s!", p.Name(), killer.Name()))
	// drop held items
	// handle RP
	rpAward := uint32(p.CurPhyStats().level() - (killer.CurPhyStats().level() - p.CurPhyStats().level()))
	plural := "point"
	if rpAward != 1 {
		plural = "points"
	}
	killer.SendMessage(fmt.Sprintf("You gain %s realm %s!", CHighlight(rpAward), plural), false)
	killer.AwardRP(rpAward)
	// create a corpse?  flag

}

func (p *Player) DamageMOB(attacker MOB, dmg uint16) {
	if (p.CurState().hits() == 0) && (dmg > 0) {
		p.KillMOB(attacker)
		return
	}
	p.CurState().adjHits(int16(-dmg), p.MaxState().hits())
	if p.CurState().hits() == 0 {
		p.SendMessage("You are almost dead!", true)
	}

	attacker.SetVictim(p)
	p.SetVictim(attacker)
}

func (p *Player) Victim() MOB {
	return p.victim
}

func (p *Player) SetVictim(newVictim MOB) {
	p.victim = newVictim
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
		p.BasePhyStats().setLevel(newLevel)
		p.RecoverPhyStats()
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

func (p *Player) releaseCorpse(w *World) {
	p.RecoverCharState()
	p.RecoverPhyStats()
	target := w.GetRoomById("A")
	if target != nil {
		p.SendMessage("You release your corpse!", true)
		p.Room().ShowOthers(p, nil, fmt.Sprintf("%s releases their corpse!", p.Name()))
		w.MoveMob(p, target)
		p.Room().ShowOthers(p, nil, fmt.Sprintf("%s appears in a puff of smoke.", p.Name()))
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
		location = p.Room().RoomID()
	} else {
		// TODO(jon): need to preserve last room
		location = ""
	}
	player := PlayerDB{
		name:       p.Name(),
		account:    p.Account,
		class:      p.BaseCharStats().CurrentClass().Name(),
		race:       p.BaseCharStats().Race().Name(),
		room:       location,
		coins:      p.Coins(),
		stre:       p.BaseCharStats().Stat(STAT_STRENGTH),
		cons:       p.BaseCharStats().Stat(STAT_CONSTITUTION),
		agil:       p.BaseCharStats().Stat(STAT_AGILITY),
		dext:       p.BaseCharStats().Stat(STAT_DEXTERITY),
		inte:       p.BaseCharStats().Stat(STAT_INTELLIGENCE),
		wisd:       p.BaseCharStats().Stat(STAT_WISDOM),
		con_loss:   0,
		level:      p.BasePhyStats().level(),
		exp:        p.Experience,
		rp:         p.RealmPoints,
		hits:       p.MaxState().hits(),
		fat:        p.MaxState().fat(),
		power:      p.MaxState().power(),
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
