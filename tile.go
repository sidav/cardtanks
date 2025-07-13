package main

type tileCode byte

const (
	TILE_FLOOR tileCode = iota
	TILE_ENEMY_SPAWNER
	TILE_ARMOR
	TILE_WALL
	TILE_DAMAGED_WALL
	TILE_FOREST
	TILE_ICE
	TILE_WATER
	TILE_FLAG
	TILE_EAGLE
)

func (tc tileCode) getSpriteAtlas() *spriteAtlas {
	switch tc {
	case TILE_WALL:
		return tileAtlaces["WALL"]
	case TILE_DAMAGED_WALL:
		return tileAtlaces["DAMAGED_WALL"]
	case TILE_ARMOR:
		return tileAtlaces["ARMORED_WALL"]
	case TILE_FOREST:
		return tileAtlaces["WOOD"]
	case TILE_ICE:
		return tileAtlaces["ICE"]
	case TILE_WATER:
		return tileAtlaces["WATER"]
	case TILE_FLAG:
		return tileAtlaces["FLAG"]
	case TILE_EAGLE:
		return tileAtlaces["EAGLE"]
	}
	return nil
}

func (tc tileCode) is(code tileCode) bool {
	return tc == code
}

func (tc tileCode) canBeDrivenOn() bool {
	return tc == TILE_FLOOR ||
		tc == TILE_ENEMY_SPAWNER ||
		tc == TILE_FOREST ||
		tc == TILE_ICE ||
		tc == TILE_WATER ||
		tc == TILE_FLAG
}

func (tc tileCode) willAiDriveOn() bool {
	return tc.canBeDrivenOn() && tc != TILE_WATER
}

func (tc tileCode) willAiShootAt() bool {
	return tc.isDestructible() && tc != TILE_EAGLE
}

func (tc tileCode) isDestructible() bool {
	return tc == TILE_WALL || tc == TILE_DAMAGED_WALL || tc == TILE_EAGLE
}

func (tc tileCode) canBeShotThrough() bool {
	return tc != TILE_WALL && tc != TILE_DAMAGED_WALL && tc != TILE_ARMOR && tc != TILE_EAGLE
}

func (tc *tileCode) destroy() {
	switch *tc {
	case TILE_WALL:
		(*tc) = TILE_DAMAGED_WALL
	case TILE_DAMAGED_WALL, TILE_FLAG, TILE_EAGLE:
		(*tc) = TILE_FLOOR
	}
}
