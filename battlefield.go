package main

import (
	"cardtanks/calc"
	"math/rand"
)

const (
	BF_WIDTH       = 12
	BF_HEIGHT      = 10
	TEAM_NONE      = 0
	TEAM_PLAYER    = 1
	TEAM_ENEMY1    = 2
	TEAM_ENEMY2    = 3
	TEAM_ENEMY3    = 4
	MAX_TEAM_CONST = 5 // needed for random team assignment
)

type battlefield struct {
	state      battlefieldState
	tiles      [BF_WIDTH][BF_HEIGHT]tile
	playerTank *tank
	tanks      []*tank

	mission             battlefieldMissionId
	missionProgress     int // general-purpose integer
	maxTanksPerTeam     int
	totalEnemyTanks     int
	spawnFastEnemies    bool
	spawnArmoredEnemies bool
}

func (b *battlefield) tileAt(x, y int) *tile {
	return &b.tiles[x][y]
}

func (b *battlefield) getTankAt(x, y int) *tank {
	if b.playerTank.x == x && b.playerTank.y == y {
		return b.playerTank
	}
	for _, t := range b.tanks {
		if t.x == x && t.y == y {
			return t
		}
	}
	return nil
}

func (b *battlefield) isPlayerTank(t *tank) bool {
	return b.playerTank == t
}

func (b *battlefield) countTilesOfType(ttype tileCode) int {
	count := 0
	for x := range b.tiles {
		for y := range b.tiles[x] {
			if b.tileAt(x, y).is(ttype) {
				count++
			}
		}
	}
	return count
}

func (b *battlefield) countTanksOfTeam(team byte) int {
	count := 0
	if b.playerTank.team == team {
		count++
	}
	for _, t := range b.tanks {
		if t.team == team {
			count++
		}
	}
	return count
}

func (b *battlefield) areTanksEnemies(t1, t2 *tank) bool {
	return t1.team != t2.team
}

func (b *battlefield) getHitCoordinatesIfTankFires(t *tank) *calc.IntVector2d {
	x, y := t.x, t.y
	for {
		x += t.dirX
		y += t.dirY
		if !b.areCoordsValid(x, y) {
			break
		}
		if b.tileAt(x, y).is(TILE_FOREST) {
			continue
		}
		tankThere := b.getTankAt(x, y)
		if tankThere != nil || !b.tileAt(x, y).canBeShotThrough() {
			vect := calc.NewIntVector2d(x, y)
			return &vect
		}
	}
	return nil
}

func (b *battlefield) areCoordsValid(x, y int) bool {
	return x >= 0 && x < len(b.tiles) && y >= 0 && y < len(b.tiles[x])
}

// Moves tank by vector, ignoring tank's direction. Returns false if the tank can't be moved.
// This is just a movement itself, it does NOT push other tanks!
func (b *battlefield) tryMovingTankByVector(t *tank, vx, vy int) bool {
	x, y := t.x+vx, t.y+vy
	if b.areCoordsValid(x, y) && b.tiles[x][y].canBeDrivenOn() && b.getTankAt(x, y) == nil {
		t.x, t.y = x, y
		t.lastMoveVector.SetEqTo(vx, vy)
	} else {
		t.lastMoveVector.SetEqTo(0, 0) // Resetting this so that the tank won't move from ice during wrong turn
		return false
	}
	return true
}

func (b *battlefield) tryPushingTankByVector(t *tank, vx, vy int) bool {
	if vx == 0 && vy == 0 {
		return false
		// panic("Vector push failed: report this")
	}
	x, y := t.x+vx, t.y+vy
	otherTank := b.getTankAt(x, y)
	if otherTank != nil { // pushing the other tank
		b.tryMovingTankByVector(otherTank, vx, vy)
	}
	return b.tryMovingTankByVector(t, vx, vy)
}

func (b *battlefield) tryPushingTankForward(t *tank) bool {
	return b.tryPushingTankByVector(t, t.dirX, t.dirY)
}

// Returns true if any tank was moved this way
func (b *battlefield) tryPushingAllTanksOnIce() bool {
	anyPushed := false
	if b.tileAt(b.playerTank.getCoords()).is(TILE_ICE) {
		pushed := b.tryPushingTankByVector(b.playerTank, b.playerTank.lastMoveVector.X, b.playerTank.lastMoveVector.Y)
		anyPushed = anyPushed || pushed
	}
	for _, t := range b.tanks {
		if b.tileAt(t.getCoords()).is(TILE_ICE) {
			pushed := b.tryPushingTankByVector(t, t.lastMoveVector.X, t.lastMoveVector.Y)
			anyPushed = anyPushed || pushed
		}
	}
	return anyPushed
}

func (b *battlefield) areAnyTanksOnIce() bool {
	if b.tileAt(b.playerTank.getCoords()).is(TILE_ICE) {
		return true
	}
	for _, t := range b.tanks {
		if b.tileAt(t.getCoords()).is(TILE_ICE) {
			return true
		}
	}
	return false
}

func (b *battlefield) lineOfFireExistsBetweenTwoTanks(t1, t2 *tank) bool {
	return b.lineOfFireExistsBetweenCoords(t1.x, t1.y, t2.x, t2.y)
}

func (b *battlefield) lineOfFireExistsBetweenCoords(x1, y1, x2, y2 int) bool {
	if !(x1 == x2 || y1 == y2) {
		return false
	}
	v := calc.NewIntVector2d(x2-x1, y2-y1)
	v.Normalize()
	x, y := x1, y1
	for {
		x += v.X
		y += v.Y
		if x == x2 && y == y2 {
			return true
		}
		tankThere := b.getTankAt(x, y)
		if !b.tiles[x][y].canBeShotThrough() || tankThere != nil {
			return false
		}
	}
}

func (b *battlefield) countTilesOfTypeAroundCoords(ttype tileCode, x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i*j != 0 || (i == 0 && j == 0) {
				continue
			}
			if b.areCoordsValid(x+i, y+j) && b.tileAt(x+i, y+j).is(ttype) {
				count++
			}
		}
	}
	return count
}

func (b *battlefield) selectRandomMapCoordsByAllowanceFunc(allowanceFunc func(x, y int) bool) *calc.IntVector2d {
	var candidates []calc.IntVector2d
	for x := range len(b.tiles) {
		for y := range len(b.tiles[x]) {
			if allowanceFunc(x, y) {
				candidates = append(candidates, calc.NewIntVector2d(x, y))
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	index := rand.Intn(len(candidates))
	return &candidates[index]
}

func (b *battlefield) trySpawnNewEnemy() bool {
	v := b.selectRandomMapCoordsByAllowanceFunc(func(x, y int) bool {
		return b.tileAt(x, y).is(TILE_FLOOR) &&
			b.getTankAt(x, y) == nil &&
			!b.lineOfFireExistsBetweenCoords(x, y, b.playerTank.x, b.playerTank.y)
	})
	if v == nil {
		return false
	}
	var newTankCode byte = TANK_ENEMY
	if b.spawnFastEnemies && rand.Intn(100) < 33 {
		newTankCode = TANK_ENEMY_FAST
	} else if b.spawnArmoredEnemies && rand.Intn(100) < 33 {
		newTankCode = TANK_ENEMY_ARMORED
	}

	newTank := createTank(newTankCode, TEAM_ENEMY1, v.X, v.Y)
	newTank.faceRandomDirection()
	b.tanks = append(b.tanks, newTank)

	b.totalEnemyTanks--
	return true
}

func (b *battlefield) trySinkingTanks() bool {
	sunken := false
	if b.tileAt(b.playerTank.getCoords()).is(TILE_WATER) {
		b.playerTank.health = 0
		sunken = true
	}
	// Handle enemies
	for _, t := range b.tanks {
		if b.tileAt(t.getCoords()).is(TILE_WATER) {
			t.health = 0
			sunken = true
		}
	}
	return sunken
}

func (b *battlefield) clearDestroyedTanks() {
	for i := len(b.tanks) - 1; i >= 0; i-- {
		t := b.tanks[i]
		if t.health <= 0 {
			b.tanks = append(b.tanks[:i], b.tanks[i+1:]...)
		}
	}
}
