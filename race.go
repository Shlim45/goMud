package main

import "fmt"

type Race interface {
	SetName(name string)
	Name() string
	SetEnabled(yesNo bool)
	Enabled() bool
	Qualifies(m *MOB) bool
	//SetRealm(realm Realm)
	//Realm() Realm
	SetStatBonuses(statBonuses [NUM_STATS]float64)
	StatBonuses() *[NUM_STATS]float64
	SaveRaceToDBQuery() string
}

type PlayerRace struct {
	name string
	//realm       Realm
	enabled     bool
	statBonuses [NUM_STATS]float64
}

func (r *PlayerRace) Name() string {
	return r.name
}

func (r *PlayerRace) SetName(name string) {
	r.name = name
}

func (r *PlayerRace) Enabled() bool {
	return r.enabled
}

func (r *PlayerRace) SetEnabled(yesNo bool) {
	r.enabled = yesNo
}

func (r *PlayerRace) Qualifies(m *MOB) bool {
	return true
}

//func (r *PlayerRace) Realm() Realm {
//	return r.realm
//}
//
//func (r *PlayerRace) SetRealm(realm Realm) {
//	r.realm = realm
//}

func (r *PlayerRace) SetStatBonuses(statBonuses [NUM_STATS]float64) {
	r.statBonuses = statBonuses
}

func (r *PlayerRace) StatBonuses() *[NUM_STATS]float64 {
	return &r.statBonuses
}

func (r *PlayerRace) SaveRaceToDBQuery() string {
	return fmt.Sprintf("INSERT INTO Race VALUES ('%s', %v) AS new ON DUPLICATE KEY UPDATE name=new.name, enabled=new.enabled",
		r.Name(), r.Enabled())
}
