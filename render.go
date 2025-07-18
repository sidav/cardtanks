package main

import (
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var DEFAULT_TINT = rl.RayWhite
var gameOverRgb color.RGBA

const TILE_SIZE_PIXELS = 16 * SPRITE_SCALE_FACTOR

type renderer struct {
	bf      *battlefield
	plr     *player
	bullets []*drawnBullet
}

func (r *renderer) renderBattlefield(b *battlefield, plr *player) {
	r.bf = b
	r.plr = plr
	if !(r.bf.state.Is(BS_SHOOT) || r.bf.state.Is(BS_PLAYER_SHOOTS_DURING_TURN)) {
		r.bullets = r.bullets[:0]
	}

	r.renderLevelOutline()
	r.renderTiles(b)

	for i := range b.tanks {
		r.renderTank(b.tanks[i])
	}
	r.renderTank(b.playerTank)

	if r.bf.state.Is(BS_SHOOT) {
		r.bf.state.Lock()
		if len(r.bullets) == 0 {
			r.createBulletsForAllTanks()
		}
		r.handleAndRenderBullets()
	}
	if r.bf.state.Is(BS_PLAYER_SHOOTS_DURING_TURN) {
		r.bf.state.Lock()
		if len(r.bullets) == 0 {
			r.createBulletForTank(b.playerTank)
		}
		r.handleAndRenderBullets()
	}

	r.renderWood(b)
	r.renderBattlefieldInfo()
	if r.bf.state.awaitsPlayerInput() {
		r.renderButtons()
		r.renderCardsInHand()
	}
}

func (r *renderer) renderTile(b *battlefield, x, y int) {
	t := b.tileAt(x, y)
	osx, osy := x*TILE_SIZE_PIXELS, y*TILE_SIZE_PIXELS
	if t.justSpawned {
		r.renderSpawnEffect(int32(osx), int32(osy))
		return
	}
	spr := t.getSpriteAtlas()
	if spr != nil {
		const changeFrameEachMs = 350
		frameNumber := int(time.Now().UnixMilli()/changeFrameEachMs) % spr.totalFrames()
		rl.DrawTexture(
			spr.getSpriteByDirectionAndFrameNumber(0, 0, frameNumber),
			int32(osx),
			int32(osy),
			DEFAULT_TINT,
		)
	}
}

func (r *renderer) renderTank(t *tank) {
	x, y := float32(t.x*TILE_SIZE_PIXELS)+1, float32(t.y*TILE_SIZE_PIXELS)+1

	if t.justSpawned {
		r.renderSpawnEffect(int32(x), int32(y))
		return
	}

	if t.health <= 0 {
		const changeFrameEachMs = 100
		frameNumber := int(time.Now().UnixMilli()/changeFrameEachMs) % effectAtlaces["EXPLOSION"].totalFrames()
		rl.DrawTexture(
			effectAtlaces["EXPLOSION"].getSpriteByDirectionAndFrameNumber(0, 0, frameNumber),
			int32(x),
			int32(y),
			rl.White,
		)
		return
	}

	tint := color.RGBA{R: 255, G: 0, B: 255, A: 255}
	switch t.team {
	case TEAM_PLAYER:
		if r.bf.isPlayerTank(t) {
			tint = color.RGBA{R: 255, G: 255, B: 64, A: 255}
		} else {

			tint = color.RGBA{R: 64, G: 255, B: 64, A: 255}
		}
	case TEAM_ENEMY1:
		tint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	case TEAM_ENEMY2:
		tint = color.RGBA{R: 255, G: 64, B: 64, A: 255}
	case TEAM_ENEMY3:
		tint = color.RGBA{R: 64, G: 64, B: 255, A: 255}
	}
	if t == r.bf.playerTank {
	} else if t.team == TEAM_PLAYER {
	}
	rl.DrawTexture(
		t.getSpriteAtlas().getSpriteByDirectionAndFrameNumber(t.dirX, t.dirY, 0),
		int32(x),
		int32(y),
		tint,
	)
}

func (r *renderer) renderTiles(b *battlefield) {
	for x := range b.tiles {
		for y := range b.tiles[x] {
			r.renderTile(b, x, y)
		}
	}
}

func (r *renderer) renderWood(b *battlefield) {
	for x := range b.tiles {
		for y, t := range b.tiles[x] {
			if t.is(TILE_FOREST) {
				r.renderTile(b, x, y)
			}
		}
	}
}

func (r *renderer) renderSpawnEffect(x, y int32) {
	const changeFrameEachMs = 100
	frameNumber := int(time.Now().UnixMilli()/changeFrameEachMs) % effectAtlaces["SPAWN"].totalFrames()
	rl.DrawTexture(
		effectAtlaces["SPAWN"].getSpriteByDirectionAndFrameNumber(0, 0, frameNumber),
		int32(x),
		int32(y),
		rl.White,
	)
}

func (r *renderer) renderLevelOutline() {
	const thickness = 5
	x, y := 0, 0
	rl.DrawRectangleLinesEx(
		rl.Rectangle{
			X:      float32(x) - thickness,
			Y:      float32(y) - thickness,
			Width:  float32(TILE_SIZE_PIXELS*len(r.bf.tiles)) + 2*thickness,
			Height: float32(TILE_SIZE_PIXELS*len(r.bf.tiles[0])) + 2*thickness,
		},
		thickness,
		color.RGBA{
			R: 64,
			G: 64,
			B: 64,
			A: 255,
		},
	)
}
