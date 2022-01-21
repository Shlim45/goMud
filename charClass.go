package main

import "fmt"

type CharClass interface {
	SetName(name string)
	Name() string
	SetEnabled(yesNo bool)
	Enabled() bool
	Qualifies(m *MOB) bool
	//AllowedRace(r *Race) bool
	SetRealm(realm Realm)
	Realm() Realm
	SetStatBonuses(statBonuses [NUM_STATS]float64)
	StatBonuses() *[NUM_STATS]float64
	SaveCharClassToDBQuery() string
}

type PlayerClass struct {
	name        string
	realm       Realm
	enabled     bool
	statBonuses [NUM_STATS]float64
	//qualifyingRaces []*Race
}

func (pc *PlayerClass) Name() string {
	return pc.name
}

func (pc *PlayerClass) SetName(name string) {
	pc.name = name
}

func (pc *PlayerClass) Enabled() bool {
	return pc.enabled
}

func (pc *PlayerClass) SetEnabled(yesNo bool) {
	pc.enabled = yesNo
}

func (pc *PlayerClass) Qualifies(m *MOB) bool {
	return true
}

func (pc *PlayerClass) Realm() Realm {
	return pc.realm
}

func (pc *PlayerClass) SetRealm(realm Realm) {
	pc.realm = realm
}

func (pc *PlayerClass) SetStatBonuses(statBonuses [NUM_STATS]float64) {
	pc.statBonuses = statBonuses
}

func (pc *PlayerClass) StatBonuses() *[NUM_STATS]float64 {
	return &pc.statBonuses
}

func (pc *PlayerClass) SaveCharClassToDBQuery() string {
	return fmt.Sprintf("INSERT INTO CharClass VALUES ('%s', %d, %v)",
		pc.Name(), uint8(pc.Realm()), pc.Enabled())
}
