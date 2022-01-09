package main

type CharStats struct {
	Strength     uint8
	Constitution uint8
	Agility      uint8
	Dexterity    uint8
	Intelligence uint8
	Wisdom       uint8
}

func (cState *CharStats) copyOf() *CharStats {
	copyOf := CharStats{
		Strength:     cState.Strength,
		Constitution: cState.Constitution,
		Agility:      cState.Agility,
		Dexterity:    cState.Dexterity,
		Intelligence: cState.Intelligence,
		Wisdom:       cState.Wisdom,
	}
	return &copyOf
}

func (cState *CharStats) strength() uint8 {
	return cState.Strength
}

func (cState *CharStats) setStrength(newStr uint8) {
	cState.Strength = newStr
}

func (cState *CharStats) constitution() uint8 {
	return cState.Constitution
}

func (cState *CharStats) setConstitution(newCon uint8) {
	cState.Constitution = newCon
}

func (cState *CharStats) agility() uint8 {
	return cState.Agility
}

func (cState *CharStats) setAgility(newAgi uint8) {
	cState.Agility = newAgi
}

func (cState *CharStats) dexterity() uint8 {
	return cState.Dexterity
}

func (cState *CharStats) setDexterity(newDex uint8) {
	cState.Dexterity = newDex
}

func (cState *CharStats) intelligence() uint8 {
	return cState.Intelligence
}

func (cState *CharStats) setIntelligence(newInt uint8) {
	cState.Intelligence = newInt
}

func (cState *CharStats) wisdom() uint8 {
	return cState.Wisdom
}

func (cState *CharStats) setWisdom(newWis uint8) {
	cState.Wisdom = newWis
}
