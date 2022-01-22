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

func (cState *CharState) adjHits(amount int16, max uint16) {
	newHits := int32(cState.hits()) + int32(amount)
	if newHits < 0 {
		cState.setHits(0)
	} else if newHits > int32(max) {
		cState.setHits(max)
	} else {
		cState.setHits(uint16(newHits))
	}
}

func (cState *CharState) adjFat(amount int16, max uint16) {
	newFat := int32(cState.fat()) + int32(amount)
	if newFat < 0 {
		cState.setFat(0)
	} else if newFat > int32(max) {
		cState.setFat(max)
	} else {
		cState.setFat(uint16(newFat))
	}
}

func (cState *CharState) adjPower(amount int16, max uint16) {
	newPower := int32(cState.power()) + int32(amount)
	if newPower < 0 {
		cState.setPower(0)
	} else if newPower > int32(max) {
		cState.setPower(max)
	} else {
		cState.setPower(uint16(newPower))
	}
}
