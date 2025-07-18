package main

import (
	"cardtanks/calc"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A bullet is ONLY A VISUAL thing here, and exists only on-screen, that's why it belongs to renderer
type drawnBullet struct {
	onScreenX, onScreenY int
	dirX, dirY           int
	targetTileV          calc.IntVector2d
}

func (r *renderer) handleAndRenderBullets() {
	allArrived := true
	for _, b := range r.bullets {
		arrived := b.targetTileV.Is(b.onScreenX/TILE_SIZE_PIXELS, b.onScreenY/TILE_SIZE_PIXELS)
		var atl *spriteAtlas
		var spr rl.Texture2D
		// Move bullet
		if !arrived {
			allArrived = false
			b.onScreenX += b.dirX * 4 * SPRITE_SCALE_FACTOR
			b.onScreenY += b.dirY * 4 * SPRITE_SCALE_FACTOR
		}
		// Draw bullets
		if arrived {
			atl = effectAtlaces["EXPLOSION"]
			spr = atl.getSpriteByDirectionAndFrameNumber(0, 0, 0)
		} else {
			atl = projectileAtlaces["BULLET"]
			spr = atl.getSpriteByDirectionAndFrameNumber(b.dirX, b.dirY, 0)
		}
		rl.DrawTexture(
			spr,
			int32(b.onScreenX-atl.spriteSize/2),
			int32(b.onScreenY-atl.spriteSize/2),
			rl.White,
		)

	}
	if allArrived {
		r.bf.state.Unlock()
	}
}

func (r *renderer) createBulletsForAllTanks() {
	if r.bf.PlayerWillTankShoot() {
		r.createBulletForTank(r.bf.playerTank)
	}
	for _, t := range r.bf.tanks {
		if r.bf.aiWillTankShoot(t) {
			r.createBulletForTank(t)
		}
	}
}

func (r *renderer) createBulletForTank(t *tank) {
	r.bullets = append(r.bullets, &drawnBullet{
		onScreenX:   t.x*TILE_SIZE_PIXELS + TILE_SIZE_PIXELS/2 - SPRITE_SCALE_FACTOR,
		onScreenY:   t.y*TILE_SIZE_PIXELS + TILE_SIZE_PIXELS/2 - SPRITE_SCALE_FACTOR,
		dirX:        t.dirX,
		dirY:        t.dirY,
		targetTileV: *r.bf.getHitCoordinatesIfTankFires(t),
	})

}
