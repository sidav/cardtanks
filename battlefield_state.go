package main

import (
	"cardtanks/calc"
	"strconv"
	"time"
)

type bfStateCode byte

const (
	BS_BEFORE_PLAYER_TURN bfStateCode = iota
	BS_PLAYER_TURN
	BS_PLAYER_MOVES
	BS_PLAYER_SHOOTS_DURING_TURN
	BS_PLAYER_ENDED_TURN
	BS_NONPLAYER_TANK_MOVES
	BS_SHOOT
	BS_SPAWN_NEW_ENEMIES
	BS_TEMP_PAUSE // Temporary state, needed just for rendering something
)

type battlefieldState struct {
	code     bfStateCode
	prevCode bfStateCode
	locked   bool // Can be used from outside game logic (e.g. renderer) to hold the state longer

	// stuff needed for the state continuation
	actionsRemaining    int
	currentEntityNumber int
	intentVector        calc.IntVector2d

	// General purpose vars, ARE RESET when switching states
	stateStartTime time.Time
	flag           bool
	msDuration     int
}

func (bs *battlefieldState) Is(code bfStateCode) bool {
	return bs.code == code
}

func (bs *battlefieldState) Lock() {
	bs.locked = true
}

func (bs *battlefieldState) Unlock() {
	bs.locked = false
}

func (bs *battlefieldState) switchTo(newState bfStateCode) {
	bs.prevCode = bs.code
	bs.code = newState
	bs.stateStartTime = time.Now()
	bs.flag = false
	bs.msDuration = 0
}

func (bs *battlefieldState) pauseFor(ms int) {
	bs.switchTo(BS_TEMP_PAUSE)
	bs.msDuration = ms
}

func (bs *battlefieldState) justUnpaused() bool {
	return bs.prevCode == BS_TEMP_PAUSE
}

func (bs *battlefieldState) resetTime() {
	bs.stateStartTime = time.Now()
}

func (bs *battlefieldState) msElapsed(ms int) bool {
	if bs.stateStartTime.IsZero() {
		bs.resetTime()
	}
	msInThisState := int(time.Since(bs.stateStartTime) / time.Millisecond)
	return msInThisState > ms
}

func (bs *battlefieldState) awaitsPlayerInput() bool {
	return bs.code == BS_PLAYER_TURN
}

func (bs *battlefieldState) currentStateName() string {
	switch bs.code {
	case BS_BEFORE_PLAYER_TURN:
		return "Cleanup"

	case BS_PLAYER_TURN:
		return "Your turn"
	case BS_NONPLAYER_TANK_MOVES:
		return "Enemy activity " + strconv.Itoa(bs.actionsRemaining)
	case BS_SHOOT:
		return "Fire"
	case BS_SPAWN_NEW_ENEMIES:
		return "Enemy reinforcements"
	case BS_TEMP_PAUSE:
		if bs.prevCode != BS_TEMP_PAUSE {
			bs.code = bs.prevCode
			name := bs.currentStateName()
			bs.code = BS_TEMP_PAUSE
			return name + "[P]"
		} else {
			return "ERROR"
		}
	default:
		return ""
	}
}
