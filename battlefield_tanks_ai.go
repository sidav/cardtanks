package main

import "cardtanks/calc"

// It's the "AI" lol
// Stupid and predictable on purpose
func (b *battlefield) actForNonplayerTank(t *tank) {
	var enemyTank *tank
	// Are we enemies to player, and if yes, do we see them?
	if b.areTanksEnemies(t, b.playerTank) && b.lineOfFireExistsBetweenTwoTanks(t, b.playerTank) {
		enemyTank = b.playerTank
	}
	// Do we see any other enemy tank?
	if enemyTank == nil {
		for _, otherTank := range b.tanks {
			if b.areTanksEnemies(t, otherTank) && b.lineOfFireExistsBetweenTwoTanks(t, otherTank) {
				enemyTank = otherTank
				break
			}
		}
	}
	// First: if an enemy can be seen...
	if enemyTank != nil {
		rotated := b.aiTryRotateTankTowardsCoords(t, enemyTank.x, enemyTank.y)
		if rotated {
			return // If rotated, current action ends here
		}
		if b.aiWillTankMoveForward(t) {
			b.tryPushingTankForward(t)
		}
		return
	}

	// Else: if the next cell is passable, move to it
	if b.aiWillTankMoveForward(t) {
		b.tryPushingTankForward(t)
	} else { // If the next cell is not passable, rotate
		t.rotateRight()
	}
}

func (b *battlefield) aiWillTankMoveForward(t *tank) bool {
	nextX, nextY := t.x+t.dirX, t.y+t.dirY
	tankThere := b.getTankAt(nextX, nextY)
	if tankThere != nil {
		return b.areTanksEnemies(t, tankThere)
	}
	return b.areCoordsValid(nextX, nextY) &&
		b.aiWillTankDriveOnTile(b.tileAt(nextX, nextY))
}

func (b *battlefield) aiWillTankShoot(t *tank) bool {
	v := b.getHitCoordinatesIfTankFires(t)
	if v == nil {
		return false
	}
	tankThere := b.getTankAt(v.X, v.Y)
	if tankThere != nil && b.areTanksEnemies(t, tankThere) {
		return true
	}
	if b.aiWillTankShootAtTile(t, b.tileAt(v.X, v.Y)) {
		return true
	}
	return false
}

func (b *battlefield) aiWillTankDriveOnTile(t *tile) bool {
	return t.canBeDrivenOn() && !t.is(TILE_WATER)
}

func  (b *battlefield) aiWillTankShootAtTile(tnk *tank, til *tile) bool {
	return til.isDestructible() && til.team != tnk.team
}

// Returns false if the tank was not rotated (e.g. already at coords)
func (b *battlefield) aiTryRotateTankTowardsCoords(t *tank, x, y int) bool {
	v := calc.NewIntVector2d(x-t.x, y-t.y)
	v.Normalize()

	if v.Is(t.dirX, t.dirY) {
		return false
	}
	// It's clunky, but well, it works
	// first, try to rotate left
	t.rotateLeft()
	if !v.Is(t.dirX, t.dirY) {
		// It didn't work, rotate right (to undo last rotation) and then rotate right again
		t.rotateRight()
		t.rotateRight()
	}
	return true
}
