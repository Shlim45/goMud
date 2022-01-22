package main

type Realm uint8

const (
	RealmImmortal = iota
	RealmEvil
	RealmChaos
	RealmGood
	RealmKaid
)

func (R Realm) String() string {
	switch R {
	case RealmImmortal:
		return "Immortal"
	case RealmEvil:
		return "Evil"
	case RealmChaos:
		return "Chaos"
	case RealmGood:
		return "Good"
	case RealmKaid:
		return "Kaid"
	default:
		return "None"
	}
}

func (R Realm) God() string {
	switch R {
	case RealmImmortal:
		return "Xyz"
	case RealmEvil:
		return "Arnak"
	case RealmChaos:
		return "Ra'Kur"
	case RealmGood:
		return "Niord"
	case RealmKaid:
		return "Abc"
	default:
		return "None"
	}
}
