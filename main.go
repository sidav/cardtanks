package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_W            = 1280
	WINDOW_H            = 720
	SPRITE_SCALE_FACTOR = 3
)

func main() {
	rl.InitWindow(WINDOW_W, WINDOW_H, "TANKS!")
	rl.SetTargetFPS(30)
	rl.SetExitKey(rl.KeyEscape)

	loadImageResources()

	g := NewGame()
	g.Run()

	rl.CloseWindow()
}
