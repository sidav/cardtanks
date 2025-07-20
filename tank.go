package main

import (
	"cardtanks/calc"
	"math/rand"
)

const (
	TANK_PLAYER = iota
	TANK_ENEMY
	TANK_ENEMY_FAST
	TANK_ENEMY_ARMORED
)

type tank struct {
	code       byte
	team       byte
	dirX, dirY int

	// This is a vector of last movement, regardless of current facing and the reason of the movement.
	// May be useful for e.g. ice logic
	lastMoveVector      calc.IntVector2d
	madeActionsThisTurn int

	x, y              int
	health            int
	additionalActions int // Makes (player actions this turn) + this value

	justSpawned bool // For the renderer; no gameplay effect
}

func createTank(code, team byte, x, y int) *tank {
	t := &tank{
		code:        code,
		team:        team,
		x:           x,
		y:           y,
		dirX:        0,
		dirY:        -1,
		health:      1,
		justSpawned: true,
	}
	switch code {
	case TANK_ENEMY_FAST:
		t.additionalActions = 1
	case TANK_ENEMY_ARMORED:
		t.health = 1
	}
	return t
}

func (t *tank) getSpriteAtlas() *spriteAtlas {
	switch t.code {
	case TANK_ENEMY:
		return tankAtlaces["TANK5"]
	case TANK_ENEMY_FAST:
		return tankAtlaces["TANK6"]
	case TANK_ENEMY_ARMORED:
		return tankAtlaces["TANK7"]
	}
	return tankAtlaces["TANK1"]
}

func (t *tank) getCoords() (int, int) {
	return t.x, t.y
}

func (t *tank) getFacing() (int, int) {
	return t.dirX, t.dirY
}

func (t *tank) getCoordsFacingAt() (int, int) {
	return t.x + t.dirX, t.y + t.dirY
}

func (t *tank) rotateLeft() {
	t.dirX, t.dirY = t.dirY, -t.dirX
}

func (t *tank) rotateRight() {
	t.dirX, t.dirY = -t.dirY, t.dirX
}

func (t *tank) faceRandomDirection() {
	times := rand.Intn(4)
	for range times {
		t.rotateRight()
	}
}
