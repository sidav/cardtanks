package main

import "math/rand"

func createBattlefieldSkirmish(spawners, teamLimit, totalEnemies int) *battlefield {
	b := &battlefield{
		maxTanksPerTeam: teamLimit,
		totalEnemyTanks: totalEnemies,
		mission:         BFM_SKIRMISH,
	}
	b.playerTank = createTank(TANK1, TEAM_PLAYER, 2+rand.Intn(6), 2+rand.Intn(6))

	// b.placeNTilesAtRandomByAllowanceFunc(spawners, TILE_ENEMY_SPAWNER, func(x, y int) bool {
	// 	return b.tileAt(x, y).code == TILE_FLOOR &&
	// 		b.countTilesOfTypeAroundCoords(TILE_ENEMY_SPAWNER, x, y) == 0 &&
	// 		b.getTankAt(x, y) == nil &&
	// 		!b.lineOfFireExistsBetweenCoords(b.playerTank.x, b.playerTank.y, x, y)
	// })

	b.placeNTilesAtRandomByAllowanceFunc(20, TILE_WALL, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(10, TILE_ARMOR, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil &&
			b.countTilesOfTypeAroundCoords(TILE_WALL, x, y) == 1 &&
			b.countTilesOfTypeAroundCoords(TILE_ARMOR, x, y) < 2 &&
			b.countTilesOfTypeAroundCoords(TILE_ENEMY_SPAWNER, x, y) == 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_FOREST, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR &&
			b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_WATER, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	b.placeNTilesAtRandomByAllowanceFunc(3, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR &&
			b.countTilesOfTypeAroundCoords(TILE_WATER, x, y) > 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(2, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	for range teamLimit {
		b.trySpawnNewEnemy()
	}
	return b
}

func createBattlefieldCaptureFlags(spawners, teamLimit int) *battlefield {
	b := &battlefield{
		maxTanksPerTeam: teamLimit,
		totalEnemyTanks: 100,
		mission:         BFM_CAPTURE_FLAGS,
	}
	b.playerTank = createTank(TANK1, TEAM_PLAYER, 2+rand.Intn(6), 2+rand.Intn(6))

	b.placeNTilesAtRandomByAllowanceFunc(20, TILE_WALL, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(10, TILE_ARMOR, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil &&
			b.countTilesOfTypeAroundCoords(TILE_WALL, x, y) == 1 &&
			b.countTilesOfTypeAroundCoords(TILE_ARMOR, x, y) < 2 &&
			b.countTilesOfTypeAroundCoords(TILE_ENEMY_SPAWNER, x, y) == 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_FOREST, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR &&
			b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_WATER, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	b.placeNTilesAtRandomByAllowanceFunc(3, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR &&
			b.countTilesOfTypeAroundCoords(TILE_WATER, x, y) > 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(2, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	for range teamLimit {
		b.trySpawnNewEnemy()
	}
	return b
}

func createBattlefieldDestroyEagles(spawners, teamLimit int) *battlefield {
	b := &battlefield{
		maxTanksPerTeam: teamLimit,
		totalEnemyTanks: 100,
		mission:         BFM_DESTROY_EAGLES,
	}
	b.playerTank = createTank(TANK1, TEAM_PLAYER, 2+rand.Intn(6), 2+rand.Intn(6))

	b.placeNTilesWithTeamAtRandomByAllowanceFunc(3, TILE_EAGLE, TEAM_ENEMY1, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR &&
			b.countTilesOfTypeAroundCoords(TILE_EAGLE, x, y) == 0 &&
			b.getTankAt(x, y) == nil &&
			!b.lineOfFireExistsBetweenCoords(b.playerTank.x, b.playerTank.y, x, y)
	})

	b.placeNTilesAtRandomByAllowanceFunc(20, TILE_WALL, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(10, TILE_ARMOR, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR && b.getTankAt(x, y) == nil &&
			b.countTilesOfTypeAroundCoords(TILE_WALL, x, y) == 1 &&
			b.countTilesOfTypeAroundCoords(TILE_ARMOR, x, y) < 2 &&
			b.countTilesOfTypeAroundCoords(TILE_ENEMY_SPAWNER, x, y) == 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_FOREST, func(x, y int) bool {
		return b.tileAt(x, y).code == TILE_FLOOR &&
			b.getTankAt(x, y) == nil
	})

	b.placeNTilesAtRandomByAllowanceFunc(5, TILE_WATER, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	b.placeNTilesAtRandomByAllowanceFunc(3, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR &&
			b.countTilesOfTypeAroundCoords(TILE_WATER, x, y) > 0
	})

	b.placeNTilesAtRandomByAllowanceFunc(2, TILE_ICE, func(x, y int) bool {
		return b.getTankAt(x, y) == nil &&
			b.tileAt(x, y).code == TILE_FLOOR
	})

	for range teamLimit {
		b.trySpawnNewEnemy()
	}
	return b
}

// Generation functions
func (b *battlefield) placeNTilesAtRandomByAllowanceFunc(count int, newCode tileCode, allowanceFunc func(x, y int) bool) {
	for range count {
		c := b.selectRandomMapCoordsByAllowanceFunc(allowanceFunc)
		if c == nil {
			return
		}
		b.tileAt(c.Unwrap()).code = newCode
	}
}

func (b *battlefield) placeNTilesWithTeamAtRandomByAllowanceFunc(count int, newCode tileCode, team byte, allowanceFunc func(x, y int) bool) {
	for range count {
		c := b.selectRandomMapCoordsByAllowanceFunc(allowanceFunc)
		if c == nil {
			return
		}
		b.tileAt(c.Unwrap()).code = newCode
		b.tileAt(c.Unwrap()).team = team
	}
}

