package main

import "slices"

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

type tile struct {
	code        tileCode
	justSpawned bool
	team        byte
}

func (t *tile) getSpriteAtlas() *spriteAtlas {
	switch t.code {
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

func (t *tile) is(code tileCode) bool {
	return t.code == code
}

func (t *tile) isOneOf(codes ...tileCode) bool {
	return slices.Contains(codes, t.code)
}

// Sets the tile to code and also sets "just spawned" status
func (t *tile) spawnAs(code tileCode) {
	t.code = code
	t.justSpawned = true
}

func (t *tile) canBeDrivenOn() bool {
	return t.code == TILE_FLOOR ||
		t.code == TILE_ENEMY_SPAWNER ||
		t.code == TILE_FOREST ||
		t.code == TILE_ICE ||
		t.code == TILE_WATER ||
		t.code == TILE_FLAG
}

func (t *tile) isDestructible() bool {
	return t.code == TILE_WALL || t.code == TILE_DAMAGED_WALL || t.code == TILE_EAGLE
}

func (t *tile) canBeShotThrough() bool {
	return t.code != TILE_WALL && t.code != TILE_DAMAGED_WALL && t.code != TILE_ARMOR && t.code != TILE_EAGLE
}

func (t *tile) destroy() {
	switch t.code {
	case TILE_WALL:
		(t.code) = TILE_DAMAGED_WALL
	case TILE_DAMAGED_WALL, TILE_FLAG, TILE_EAGLE:
		(t.code) = TILE_FLOOR
	}
}
