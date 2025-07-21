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
	BFM_BOSS_FIGHT
	BFM_MISSIONS_COUNT
)

const (
	MISSION_FLAGS_TO_CAPTURE  = 3
	MISSION_EAGLES_TO_DESTROY = 3
)

func (b *battlefield) GetMissionBriefing() (missTitle, missDescr string) {
	switch b.mission {
	case BFM_SKIRMISH:
		missTitle = "Skirmish"
		missDescr = fmt.Sprintf(
			"Destroy %d enemy tanks to win\n"+
				"Max %d enemies will appear at the same time",
			b.totalEnemyTanks+b.countTanksOfTeam(TEAM_ENEMY1), b.maxTanksPerTeam)
	case BFM_CAPTURE_FLAGS:
		missTitle = "Race for the flags"
		missDescr = fmt.Sprintf(
			"Capture %d flags to win\n"+
				"Max %d enemies will appear at the same time",
			MISSION_FLAGS_TO_CAPTURE, b.maxTanksPerTeam)
	case BFM_DESTROY_EAGLES:
		missTitle = "One-tank army"
		missDescr = fmt.Sprintf(
			"Destroy %d enemy eagles to win\n"+
				"Max %d enemies will appear at the same time, and +1 for each destroyed eagle",
			MISSION_EAGLES_TO_DESTROY, b.maxTanksPerTeam)
	case BFM_BOSS_FIGHT:
		missTitle = "Overlord"
		missDescr = fmt.Sprintf(
			"Destroy Big Bad Boss to win\n"+
				"Max %d enemies will appear at the same time, and +1 for each time the boss is hit",
			b.maxTanksPerTeam)
	}
	if b.spawnFastEnemies {
		missDescr += "\n\n Fast tanks may appear!"
	}
	if b.spawnArmoredEnemies {
		missDescr += "\n\n Armored tanks may appear!"
	}
	return
}

func (b *battlefield) GetMissionProgressString() string {
	switch b.mission {
	case BFM_SKIRMISH:
		return fmt.Sprintf(
			"%d tanks left to destroy", b.totalEnemyTanks+b.countTanksOfTeam(TEAM_ENEMY1))
	case BFM_CAPTURE_FLAGS:
		return fmt.Sprintf("%d/%d flags taken", b.missionProgress, MISSION_FLAGS_TO_CAPTURE)
	case BFM_DESTROY_EAGLES:
		return fmt.Sprintf("%d eagles left to destroy", b.countTilesOfType(TILE_EAGLE))
	case BFM_BOSS_FIGHT:
		return fmt.Sprintf("Boss health: %d/%d", b.enemyBossTank.health, b.enemyBossTank.GetMaxHealth())
	}
	return "No description"
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
		b.missionProgress = MISSION_EAGLES_TO_DESTROY - b.countTilesOfType(TILE_EAGLE)
		b.maxTanksPerTeam += b.missionProgress
	case BFM_BOSS_FIGHT:
		b.maxTanksPerTeam -= b.missionProgress
		b.missionProgress = b.enemyBossTank.GetMaxHealth() - b.enemyBossTank.health
		b.maxTanksPerTeam += b.missionProgress
	}
}

func (b *battlefield) IsMissionWon() bool {
	switch b.mission {
	case BFM_SKIRMISH:
		return b.totalEnemyTanks+b.countTanksOfTeam(TEAM_ENEMY1)+b.countTanksOfTeam(TEAM_ENEMY2)+b.countTanksOfTeam(TEAM_ENEMY3) == 0
	case BFM_CAPTURE_FLAGS:
		return b.missionProgress == MISSION_FLAGS_TO_CAPTURE
	case BFM_DESTROY_EAGLES:
		return b.countTilesOfType(TILE_EAGLE) == 0
	case BFM_BOSS_FIGHT:
		return b.enemyBossTank.health <= 0
	}
	return false
}

func (b *battlefield) IsMissionLost() bool {
	if b.playerTank.health <= 0 {
		return true
	}
	return false
}
