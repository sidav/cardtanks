package main

func (b *battlefield) actOnState(plr *player) {
	b.doMissionSpecificCheck()
	st := &b.state
	switch st.code {
	case BS_BEFORE_PLAYER_TURN:
		if !b.state.msElapsed(500) {
			return
		}
		b.doBeginningOfTurnCleanup(plr)
		st.switchTo(BS_PLAYER_TURN)

	case BS_PLAYER_TURN:
		if !st.justUnpaused() {
			sunken := b.trySinkingTanks()
			if sunken {
				st.pauseFor(250)
				return
			}
		}
		b.clearDestroyedTanks()
		b.cleanJustSpawnedStatuses()
		return

	case BS_PLAYER_MOVES:
		if b.areAnyTanksOnIce() {
			anythingPushed := b.tryPushingAllTanksOnIce()
			if anythingPushed {
				st.pauseFor(250)
				return
			}
		}
		if st.actionsRemaining == 0 {
			st.switchTo(BS_PLAYER_TURN)
			return
		}
		if st.msElapsed(200) {
			b.clearDestroyedTanks()
			moved := b.tryPushingTankByVector(b.playerTank, st.intentVector.X, st.intentVector.Y)
			st.actionsRemaining--
			st.resetTime()
			if moved {
				st.pauseFor(200)
			}
		}
	case BS_PLAYER_SHOOTS_DURING_TURN:
		if st.locked {
			st.resetTime()
		}
		if !st.msElapsed(500) {
			return
		}
		// Do the shooting
		b.doShootingForTank(b.playerTank, plr)
		b.clearDestroyedTanks()
		st.switchTo(BS_PLAYER_TURN)

	case BS_PLAYER_ENDED_TURN:
		st.currentEntityNumber = 0
		st.actionsRemaining = plr.actionsSpentForTurn
		st.switchTo(BS_NONPLAYER_TANK_MOVES)

	case BS_NONPLAYER_TANK_MOVES:
		if plr.actionsSpentForTurn == 0 {
			st.switchTo(BS_SHOOT)
			return
		}
		b.clearDestroyedTanks()
		if b.areAnyTanksOnIce() {
			anythingPushed := b.tryPushingAllTanksOnIce()
			if anythingPushed {
				st.pauseFor(200)
				return
			}
		}
		if st.actionsRemaining == 0 {
			st.currentEntityNumber++
			st.actionsRemaining = plr.actionsSpentForTurn
		}
		if st.currentEntityNumber >= len(b.tanks) {
			if st.msElapsed(300) {
				st.switchTo(BS_SHOOT)
			}
			return
		}
		b.actForNonplayerTank(b.tanks[st.currentEntityNumber])
		st.actionsRemaining--
		st.resetTime()
		b.trySinkingTanks()
		st.pauseFor(250)

	case BS_SHOOT:
		if st.locked {
			st.resetTime()
		}
		if !st.msElapsed(500) {
			return
		}
		// Do the shooting
		b.doShootingForTank(b.playerTank, plr)
		for _, t := range b.tanks {
			b.doShootingForTank(t, plr)
		}
		b.clearDestroyedTanks()
		st.switchTo(BS_SPAWN_NEW_ENEMIES)

	case BS_SPAWN_NEW_ENEMIES:
		enemiesToSpawn := b.maxTanksPerTeam - b.countTanksOfTeam(TEAM_ENEMY1)
		enemiesToSpawn = min(enemiesToSpawn, b.totalEnemyTanks)
		for range enemiesToSpawn {
			b.trySpawnNewEnemy()
		}
		st.switchTo(BS_BEFORE_PLAYER_TURN)

	///
	/// STATES WITH RETURN
	///
	// Non-gameplay state, needed for rendering pause
	case BS_TEMP_PAUSE:
		if st.msElapsed(st.msDuration) {
			st.switchTo(st.prevCode)
		}
	}
}

func (b *battlefield) PlayerWillTankShoot() bool {
	v := b.getHitCoordinatesIfTankFires(b.playerTank)
	if v == nil {
		return false
	}
	tankThere := b.getTankAt(v.X, v.Y)
	if tankThere != nil && !b.areTanksEnemies(b.playerTank, tankThere) {
		return false
	}
	return b.tileAt(v.Unwrap()).team != b.playerTank.team
}

func (b *battlefield) doShootingForTank(t *tank, plr *player) {
	if b.isPlayerTank(t) {
		if !b.PlayerWillTankShoot() {
			return
		}

	} else if !b.aiWillTankShoot(t) {
		return
	}
	v := b.getHitCoordinatesIfTankFires(t)
	if v == nil {
		return
	}
	x, y := v.Unwrap()
	hitTank := b.getTankAt(x, y)
	if hitTank != nil {
		if b.playerTank == hitTank {
			b.handlePlayerBeingHit(plr)
		} else {
			hitTank.health--
		}
		return
	}
	b.tiles[x][y].destroy()
}

func (b *battlefield) doBeginningOfTurnCleanup(plr *player) {
	plr.actionsSpentForTurn = 0
	b.cleanJustSpawnedStatuses()

}

func (b *battlefield) cleanJustSpawnedStatuses() {
	for x := range b.tiles {
		for y := range b.tiles[x] {
			b.tileAt(x, y).justSpawned = false
		}
	}
	b.playerTank.justSpawned = false
	for _, t := range b.tanks {
		t.justSpawned = false
	}
}
