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
	Stats [NUM_STATS]uint8
}

func (cState *CharStats) copyOf() *CharStats {
	var statsCopy [NUM_STATS]uint8
	for stat, value := range cState.Stats {
		statsCopy[stat] = value
	}
	copyOf := CharStats{
		Stats: statsCopy,
	}
	return &copyOf
}

func (cState *CharStats) AllStats() *[NUM_STATS]uint8 {
	return &cState.Stats
}

func (cState *CharStats) strength() uint8 {
	return cState.Stats[STAT_STRENGTH]
}

func (cState *CharStats) setStrength(newStr uint8) {
	cState.Stats[STAT_STRENGTH] = newStr
}

func (cState *CharStats) constitution() uint8 {
	return cState.Stats[STAT_CONSTITUTION]
}

func (cState *CharStats) setConstitution(newCon uint8) {
	cState.Stats[STAT_CONSTITUTION] = newCon
}

func (cState *CharStats) agility() uint8 {
	return cState.Stats[STAT_AGILITY]
}

func (cState *CharStats) setAgility(newAgi uint8) {
	cState.Stats[STAT_AGILITY] = newAgi
}

func (cState *CharStats) dexterity() uint8 {
	return cState.Stats[STAT_DEXTERITY]
}

func (cState *CharStats) setDexterity(newDex uint8) {
	cState.Stats[STAT_DEXTERITY] = newDex
}

func (cState *CharStats) intelligence() uint8 {
	return cState.Stats[STAT_INTELLIGENCE]
}

func (cState *CharStats) setIntelligence(newInt uint8) {
	cState.Stats[STAT_INTELLIGENCE] = newInt
}

func (cState *CharStats) wisdom() uint8 {
	return cState.Stats[STAT_WISDOM]
}

func (cState *CharStats) setWisdom(newWis uint8) {
	cState.Stats[STAT_WISDOM] = newWis
}
