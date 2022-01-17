package main

import (
	"log"
	"math"
)

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

func (cState *CharStats) Stat(stat Stat) uint8 {
	return cState.stats[stat]
}

func (cState *CharStats) StatBonus(stat uint8) uint8 {
	if stat >= NUM_STATS {
		return 0
	}
	bonus := float64(cState.Stats()[stat]) * cState.charClass.StatBonuses()[stat]
	return uint8(math.Round(bonus))
}
