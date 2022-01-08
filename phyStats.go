package main

type PhyStats struct {
	Level        uint8
	Attack       uint16
	Damage       uint16
	Evasion      uint16
	Defense      uint16
	MagicAttack  uint16
	MagicDamage  uint16
	MagicEvasion uint16
	MagicDefense uint16
}

func (pStats *PhyStats) level() uint8 {
	return pStats.Level
}

func (pStats *PhyStats) setLevel(newLevel uint8) {
	pStats.Level = newLevel
}

func (pStats *PhyStats) attack() uint16 {
	return pStats.Attack
}

func (pStats *PhyStats) setAttack(newAttack uint16) {
	pStats.Attack = newAttack
}

func (pStats *PhyStats) damage() uint16 {
	return pStats.Damage
}

func (pStats *PhyStats) setDamage(newDamage uint16) {
	pStats.Damage = newDamage
}

func (pStats *PhyStats) evasion() uint16 {
	return pStats.Evasion
}

func (pStats *PhyStats) setEvasion(newEvasion uint16) {
	pStats.Evasion = newEvasion
}

func (pStats *PhyStats) defense() uint16 {
	return pStats.Defense
}

func (pStats *PhyStats) setDefense(newDefense uint16) {
	pStats.Defense = newDefense
}

func (pStats *PhyStats) magicAttack() uint16 {
	return pStats.MagicAttack
}

func (pStats *PhyStats) setMagicAttack(newAttack uint16) {
	pStats.MagicAttack = newAttack
}

func (pStats *PhyStats) magicDamage() uint16 {
	return pStats.MagicDamage
}

func (pStats *PhyStats) setMagicDamage(newDamage uint16) {
	pStats.MagicDamage = newDamage
}

func (pStats *PhyStats) magicEvasion() uint16 {
	return pStats.MagicEvasion
}

func (pStats *PhyStats) setMagicEvasion(newEvasion uint16) {
	pStats.MagicEvasion = newEvasion
}

func (pStats *PhyStats) magicDefense() uint16 {
	return pStats.MagicDefense
}

func (pStats *PhyStats) setMagicDefense(newDefense uint16) {
	pStats.MagicDefense = newDefense
}

func (pStats *PhyStats) copyOf() *PhyStats {
	copyOf := PhyStats{
		Level:        pStats.Level,
		Attack:       pStats.Attack,
		Damage:       pStats.Damage,
		Evasion:      pStats.Evasion,
		Defense:      pStats.Defense,
		MagicAttack:  pStats.MagicAttack,
		MagicDamage:  pStats.MagicDamage,
		MagicEvasion: pStats.MagicEvasion,
		MagicDefense: pStats.MagicDefense,
	}
	return &copyOf
}
