package main

import (
	"cardtanks/calc"
	"math/rand"
)

const (
	TANK1 = iota
)

type tank struct {
	code       byte
	team       byte
	dirX, dirY int

	// This is a vector of last movement, regardless of current facing and the reason of the movement.
	// May be useful for e.g. ice logic
	lastMoveVector calc.IntVector2d

	x, y   int
	health int

	justSpawned bool // For the renderer; no gameplay effect
}

func createTank(code, team byte, x, y int) *tank {
	t := &tank{
		code: code,
		team: team,
		x:    x,
		y:    y,
		dirX: 0,
		dirY: -1,
	}
	t.health = 1
	t.justSpawned = true
	return t
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

func (t *tank) getSpriteAtlas() *spriteAtlas {
	return tankAtlaces["TANK1"]
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

