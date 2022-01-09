package main

type CharState struct {
	Hits     uint16
	Fat      uint16
	Power    uint16
	Alive    bool
	Standing bool
	Sitting  bool
	Laying   bool
}

func (cState *CharState) copyOf() *CharState {
	copyOf := CharState{
		Hits:     cState.Hits,
		Fat:      cState.Fat,
		Power:    cState.Power,
		Alive:    cState.Alive,
		Standing: cState.Standing,
		Sitting:  cState.Sitting,
		Laying:   cState.Laying,
	}
	return &copyOf
}

func (cState *CharState) hits() uint16 {
	return cState.Hits
}

func (cState *CharState) setHits(newHits uint16) {
	cState.Hits = newHits
}

func (cState *CharState) fat() uint16 {
	return cState.Fat
}

func (cState *CharState) setFat(newFat uint16) {
	cState.Fat = newFat
}

func (cState *CharState) power() uint16 {
	return cState.Power
}

func (cState *CharState) setPower(newPower uint16) {
	cState.Power = newPower
}

func (cState *CharState) alive() bool {
	return cState.Alive
}

func (cState *CharState) setAlive(isAlive bool) {
	cState.Alive = isAlive
}

func (cState *CharState) standing() bool {
	return cState.Standing
}

func (cState *CharState) setStanding(isStanding bool) {
	cState.Standing = isStanding
}

func (cState *CharState) sitting() bool {
	return cState.Sitting
}

func (cState *CharState) setSitting(isSitting bool) {
	cState.Sitting = isSitting
}

func (cState *CharState) laying() bool {
	return cState.Laying
}

func (cState *CharState) setLaying(isLaying bool) {
	cState.Laying = isLaying
}
