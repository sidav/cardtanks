package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/draw"
	"image/png"
	"os"
)

var (
	defaultFont rl.Font
	tankAtlaces       = map[string]*spriteAtlas{}
	tileAtlaces       = map[string]*spriteAtlas{}
	projectileAtlaces = map[string]*spriteAtlas{}
	effectAtlaces     = map[string]*spriteAtlas{}
	bonusAtlaces      = map[int]*spriteAtlas{}
)

func loadImageResources() {
	defaultFont = rl.LoadFont("assets/flexi.ttf")
	var leftXForTank = 0
	tankAtlaces["TANK1"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*0, 16, 16, 2, true)
	tankAtlaces["TANK2"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*1, 16, 16, 2, true)
	tankAtlaces["TANK3"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*2, 16, 16, 2, true)
	tankAtlaces["TANK4"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*3, 16, 16, 2, true)
	tankAtlaces["TANK5"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*4, 16, 16, 2, true)
	tankAtlaces["TANK6"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*5, 16, 16, 2, true)
	tankAtlaces["TANK7"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*6, 16, 16, 2, true)
	tankAtlaces["TANK8"] = CreateAtlasFromFile("assets/tanks.png", leftXForTank, 16*7, 16, 16, 2, true)

	tileAtlaces["WALL"] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*0, 16, 16, 1, false)
	tileAtlaces["DAMAGED_WALL"] = CreateAtlasFromFile("assets/sprites.png", 16*3, 16*0, 16, 16, 1, false)
	tileAtlaces["ARMORED_WALL"] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*1, 16, 16, 1, false)
	tileAtlaces["WATER"] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*3, 16, 16, 2, false)
	tileAtlaces["WOOD"] = CreateAtlasFromFile("assets/sprites.png", 16*1, 16*2, 16, 16, 1, false)
	tileAtlaces["ICE"] = CreateAtlasFromFile("assets/sprites.png", 16*2, 16*2, 16, 16, 1, false)
	tileAtlaces["EAGLE"] = CreateAtlasFromFile("assets/sprites.png", 16*3, 16*2, 16, 16, 1, false)
	tileAtlaces["FLAG"] = CreateAtlasFromFile("assets/sprites.png", 16*4, 16*2, 16, 16, 1, false)

	// bonusAtlaces[BONUS_HELM] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_CLOCK] = CreateAtlasFromFile("assets/sprites.png", 16*1, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_SHOVEL] = CreateAtlasFromFile("assets/sprites.png", 16*2, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_STAR] = CreateAtlasFromFile("assets/sprites.png", 16*3, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_GRENADE]= CreateAtlasFromFile("assets/sprites.png", 16*4, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_TANK]= CreateAtlasFromFile("assets/sprites.png", 16*5, 16*5, 16, 16, 1, false)
	// bonusAtlaces[BONUS_GUN]= CreateAtlasFromFile("assets/sprites.png", 16*6, 16*5, 16, 16, 1, false)
	//
	projectileAtlaces["BULLET"] = CreateAtlasFromFile("assets/projectiles.png", 0, 0, 8, 8, 2, true)
	// projectileAtlaces[PROJ_ROCKET] = CreateAtlasFromFile("assets/projectiles.png", 0, 8, 8, 8, 2, true)
	// projectileAtlaces[PROJ_LIGHTNING] = CreateAtlasFromFile("assets/projectiles.png", 0, 16, 8, 8, 2, true)
	// projectileAtlaces[PROJ_BIG] = CreateAtlasFromFile("assets/projectiles.png", 0, 24, 8, 8, 2, true)
	// projectileAtlaces[PROJ_ANNIHILATOR] = CreateAtlasFromFile("assets/projectiles.png", 0, 32, 8, 8, 2, true)
	//
	effectAtlaces["EXPLOSION"] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*6, 16, 16, 3, false)
	// effectAtlaces[EFFECT_BIG_EXPLOSION] = CreateAtlasFromFile("assets/sprites.png", 16*3, 16*6, 32, 32, 2, false)
	effectAtlaces["SPAWN"] = CreateAtlasFromFile("assets/sprites.png", 16*0, 16*4, 16, 16, 4, false)
}

//func unloadResources() {
//	for k, v := range tankAtlaces {
//		debugWritef("Unload: ", k)
//		rl.UnloadTexture(v.atlas)
//	}
//}

func extractSubimageFromImage(img image.Image, fromx, fromy, w, h int) image.Image {
	minx, miny := img.Bounds().Min.X, img.Bounds().Min.Y
	//maxx, maxy := img.Bounds().Min.X, img.Bounds().Max.Y
	subImg := img.(*image.NRGBA).SubImage(
		image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
	)
	// reset img bounds, because RayLib goes nuts about it otherwise
	subImg.(*image.NRGBA).Rect = image.Rect(0, 0, w, h)
	return subImg
}

func CreateAtlasFromFile(filename string, topleftx, toplefty, originalSpriteSize, desiredSpriteSize, totalFrames int, createAllDirections bool) *spriteAtlas {

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	newAtlas := spriteAtlas{
		spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
	}
	if createAllDirections {
		newAtlas.atlas = make([][]rl.Texture2D, 4)
	} else {
		newAtlas.atlas = make([][]rl.Texture2D, 1)
	}
	// newAtlas.atlas
	for currFrame := 0; currFrame < totalFrames; currFrame++ {
		currPic := extractSubimageFromImage(img, topleftx+currFrame*originalSpriteSize, toplefty, originalSpriteSize, originalSpriteSize)
		rlImg := rl.NewImageFromImage(currPic)
		rl.ImageResizeNN(rlImg, int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR))
		newAtlas.atlas[0] = append(newAtlas.atlas[0], rl.LoadTextureFromImage(rlImg))
		if createAllDirections {
			for i := 1; i < 4; i++ {
				rl.ImageRotateCCW(rlImg)
				newAtlas.atlas[i] = append(newAtlas.atlas[i], rl.LoadTextureFromImage(rlImg))
			}
		}
	}

	return &newAtlas
}

func textureToGolangImage(t rl.Texture2D) *image.RGBA {
	img := rl.LoadImageFromTexture(t)
	nrgba := img.ToImage().(*image.RGBA)
	return nrgba
}

func mergeImages(newImg, legs, bodies, guns image.Image, partSize int) {
	newImg.(*image.RGBA).Rect = image.Rect(0, 0, partSize, partSize)
	draw.Draw(newImg.(*image.RGBA), image.Rect(0, 0, partSize, partSize), legs, image.Point{0, 0}, draw.Over)
	newImg.(*image.RGBA).Rect = image.Rect(0, 0, partSize, partSize)
	draw.Draw(newImg.(*image.RGBA), image.Rect(0, 0, partSize, partSize), bodies, image.Point{0, 0}, draw.Over)
	newImg.(*image.RGBA).Rect = image.Rect(0, 0, partSize, partSize)
	draw.Draw(newImg.(*image.RGBA), image.Rect(0, 0, partSize, partSize), guns, image.Point{0, 0}, draw.Over)
	newImg.(*image.RGBA).Rect = image.Rect(0, 0, partSize, partSize)
}
