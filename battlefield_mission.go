package main

import (
	"cardtanks/calc"
	"fmt"
)

type battlefieldMissionId byte

const (
	BFM_SKIRMISH battlefieldMissionId = iota
	BFM_CAPTURE_FLAGS
	BFM_DESTROY_EAGLES
	BFM_MISSIONS_COUNT
)

func (b *battlefield) GetMissionDescription() (missTitle, missDescr string) {
	switch b.mission {
	case BFM_SKIRMISH:
		missTitle = "Skirmish"
		missDescr = fmt.Sprintf("Destroy %d enemy tanks to win", b.totalEnemyTanks+b.countTanksOfTeam(TEAM_ENEMY1))
	case BFM_CAPTURE_FLAGS:
		missTitle = "Race for the flags"
		missDescr = "Capture 3 flags to win"
	case BFM_DESTROY_EAGLES:
		missTitle = "One-tank army"
		missDescr = "Destroy 3 enemy bases to win. Enemy receives more reinforcements after each base gets destroyed!"
	}
	return
}

func (b *battlefield) doMissionSpecificCheck() {
	switch b.mission {
	case BFM_CAPTURE_FLAGS:
		// Check if the player picks up the flag
		if b.tileAt(b.playerTank.getCoords()).is(TILE_FLAG) {
			b.tiles[b.playerTank.x][b.playerTank.y].destroy()
			b.missionProgress++
		}
		// Spawn a flag if no flags are present
		if b.countTilesOfType(TILE_FLAG) == 0 {
			b.placeNTilesAtRandomByAllowanceFunc(1, TILE_FLAG, func(x, y int) bool {
				return b.tileAt(x, y).is(TILE_FLOOR) && calc.ApproxDistanceInt(b.playerTank.x, b.playerTank.y, x, y) >= 3
			})
		}
	case BFM_DESTROY_EAGLES:
		b.maxTanksPerTeam -= b.missionProgress
		b.missionProgress = 3 - b.countTilesOfType(TILE_EAGLE)
		b.maxTanksPerTeam += b.missionProgress
	}
}

func (b *battlefield) IsMissionWon() bool {
	switch b.mission {
	case BFM_SKIRMISH:
		return b.totalEnemyTanks+b.countTanksOfTeam(TEAM_ENEMY1)+b.countTanksOfTeam(TEAM_ENEMY2)+b.countTanksOfTeam(TEAM_ENEMY3) == 0
	case BFM_CAPTURE_FLAGS:
		return b.missionProgress == 3
	case BFM_DESTROY_EAGLES:
		return b.countTilesOfType(TILE_EAGLE) == 0
	}
	return false
}

func (b *battlefield) IsMissionLost() bool {
	if b.playerTank.health <= 0 {
		return true
	}
	return false
}
