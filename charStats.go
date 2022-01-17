package main

import "log"

type Stat uint8

const (
	STAT_STRENGTH = iota
	STAT_CONSTITUTION
	STAT_AGILITY
	STAT_DEXTERITY
	STAT_INTELLIGENCE
	STAT_WISDOM
	NUM_STATS
)

func StatToString(S uint8) string {
	switch S {
	case STAT_STRENGTH:
		return "Strength"
	case STAT_CONSTITUTION:
		return "Constitution"
	case STAT_AGILITY:
		return "Agility"
	case STAT_DEXTERITY:
		return "Dexterity"
	case STAT_INTELLIGENCE:
		return "Intelligence"
	case STAT_WISDOM:
		return "Wisdom"
	default:
		log.Fatalln("ERROR: Stat.String() default case")
		return ""
	}
}

type CharStats struct {
	stats     [NUM_STATS]uint8
	charClass CharClass
}

func (cState *CharStats) CurrentClass() CharClass {
	return cState.charClass
}

func (cState *CharStats) copyOf() *CharStats {
	var statsCopy [NUM_STATS]uint8
	for stat, value := range cState.stats {
		statsCopy[stat] = value
	}
	copyOf := CharStats{
		stats:     statsCopy,
		charClass: cState.CurrentClass(),
	}
	return &copyOf
}

func (cState *CharStats) SetStat(stat uint8, value uint8) {
	if stat < NUM_STATS {
		cState.stats[stat] = value
	}
}

func (cState *CharStats) Stats() *[NUM_STATS]uint8 {
	return &cState.stats
}

//
//func (cState *CharStats) strength() uint8 {
//	return cState.stats[STAT_STRENGTH]
//}
//
//func (cState *CharStats) setStrength(newStr uint8) {
//	cState.stats[STAT_STRENGTH] = newStr
//}
//
//func (cState *CharStats) constitution() uint8 {
//	return cState.stats[STAT_CONSTITUTION]
//}
//
//func (cState *CharStats) setConstitution(newCon uint8) {
//	cState.stats[STAT_CONSTITUTION] = newCon
//}
//
//func (cState *CharStats) agility() uint8 {
//	return cState.stats[STAT_AGILITY]
//}
//
//func (cState *CharStats) setAgility(newAgi uint8) {
//	cState.stats[STAT_AGILITY] = newAgi
//}
//
//func (cState *CharStats) dexterity() uint8 {
//	return cState.stats[STAT_DEXTERITY]
//}
//
//func (cState *CharStats) setDexterity(newDex uint8) {
//	cState.stats[STAT_DEXTERITY] = newDex
//}
//
//func (cState *CharStats) intelligence() uint8 {
//	return cState.stats[STAT_INTELLIGENCE]
//}
//
//func (cState *CharStats) setIntelligence(newInt uint8) {
//	cState.stats[STAT_INTELLIGENCE] = newInt
//}
//
//func (cState *CharStats) wisdom() uint8 {
//	return cState.stats[STAT_WISDOM]
//}
//
//func (cState *CharStats) setWisdom(newWis uint8) {
//	cState.stats[STAT_WISDOM] = newWis
//}
