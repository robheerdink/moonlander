package src

import (
	//"image/color"
	"bytes"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	ass "moonlander/assets"
	com "moonlander/src/components"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/lafriks/go-tiled"
)

// Variables related to level
var (
	background     *ebiten.Image
	spaceShip      *ebiten.Image
	player         com.Player
	DrawWorldList  []com.Drawer
	DrawScreenList []com.Drawer
	HitAbleList    []com.HitAble
	UpdateList     []com.Updater
	CollideList    []com.Collider
)

const mapPath = "assets/tiled/tiled-squares.tmx" // Path to your Tiled Map.

// ClearLevel clears are references from level
func ClearLevel() {
	background = nil
	spaceShip = nil
	DrawWorldList = nil
	DrawScreenList = nil
	HitAbleList = nil
	UpdateList = nil
	CollideList = nil
}

func mapIndex(x, y int, m *tiled.Map) int {
	return x + y*m.Width
}

func mapXY(index int, m *tiled.Map) (int, int) {
	x := index % m.Width
	y := int(index / m.Width)
	return x, y
}

// LoadTMX loads a level with format TMX (tiled)
func LoadTMX(mapPath string) {

	// Parse .tmx file.
	m, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	// loop through tile layers
	for _, layer := range m.Layers {
		fmt.Println(layer.Name)
		for i, tile := range layer.Tiles {
			//tx, ty := mapXY(i, m)
			x, y := layer.GetTilePosition(i)
			if tile.Tileset != nil {
				levelItem(strconv.FormatUint(uint64(tile.ID), 10), "", x, y, m.TileWidth, m.TileHeight, 0, tile.Tileset.Tiles[tile.ID].Properties)
			}
		}
	}
	// loop through object layers
	for _, objLayer := range m.ObjectGroups {
		fmt.Println(objLayer.Name)
		for _, obj := range objLayer.Objects {
			levelItem(obj.Type, obj.Name, int(obj.X), int(obj.Y), int(obj.Width), int(obj.Height), int(obj.Rotation), obj.Properties)
		}
	}
}

func levelItem(id, name string, x, y, w, h, rotation int, prop tiled.Properties) {
	/* id   = object.type or tile.id
	   name = object.name or nil
	*/
	fmt.Println(id, name, x, y, w, h, rotation, prop)

}

// LoadLevel loads a specific level
func LoadLevel(name string) {

	LoadTMX("assets/tiled/tiled-squares.tmx")

	// common for all levels
	preloadImages()

	sha.LP = sha.LevelProperties{
		Gravity:    0.000,
		Friction:   0.992,
		Width:      1150,
		Height:     864,
		Lives:      3,
		CurrentLap: 0,
		MaxLaps:    3,
		LapTimes:   []time.Duration{},
	}

	// some abbrivations
	V := com.Vector{}
	var size int = 16
	var cpSize int = 4

	//bg := com.NewSprite(background, 0, 0, 0, V)
	tb := com.NewTextBlock(sha.LP.Width-150, 50)
	px, py := 0, 0

	img, _, err := image.Decode(bytes.NewReader(ass.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	if name == "lvl01" {
		sha.LP.Gravity = 0.015
		sha.LP.Width = sha.ScreenWidth * 2
		sha.LP.Height = sha.ScreenHeight * 2
		LW, LH, HLW, HLH := getLevelDimensionsInt()
		SetWorldImage(LW, LH)

		bg := com.NewBackground(background, 0, 0, 0, LW, LH, V)
		runnerImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		anim := com.NewAnim(runnerImage, 200, 200, 0, V, com.NewFrame(0, 32, 32, 32, 8, 5))
		px, py = HLW-300, HLH+100
		player = com.NewPlayer(sha.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, sha.Red50)
		cp1 := com.NewCheckpoint(sha.IDCheckpoint, HLW-cpSize, HLH-cpSize/2, HLW, cpSize, sha.Cyan25, true)
		cp2 := com.NewCheckpoint(sha.IDCheckpoint, HLW-cpSize/2, 0, cpSize, HLH, sha.Green25, true)
		cp3 := com.NewCheckpoint(sha.IDCheckpoint, 0, HLH-cpSize/2, HLW, cpSize, sha.Yellow25, true)
		finish := com.NewFinish(sha.IDFinish, HLW-cpSize/2, HLH, cpSize, HLH, sha.White25, []*com.Checkpoint{&cp1, &cp2, &cp3})
		wallT := com.NewWall(sha.IDWall, 0, 0, LW-size, size, sha.Blue50)
		wallL := com.NewWall(sha.IDWall, 0, size, size, LH-size, sha.Blue50)
		wallB := com.NewWall(sha.IDWall, size, LH-size, LW-size, size, sha.Blue50)
		wallR := com.NewWall(sha.IDWall, LW-size, 0, size, LH-size, sha.Blue50)
		wallCL := com.NewWall(sha.IDWall, 300, HLH-size/2, HLW-300, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-300, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 300-size/2, size, HLH-300, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-300, sha.Purple50)
		// add to interface lists
		DrawWorldList = append(DrawWorldList, &bg, &anim, &player, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		DrawScreenList = append(DrawScreenList, &tb)
		HitAbleList = append(HitAbleList, &player, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &anim)
		CollideList = append(CollideList, &player)
		//spwanRandomSquares(HitAbleList, 2, 50)

	} else if name == "lvl02" {
		sha.LP.Gravity = 0.015
		sha.LP.Width = sha.ScreenWidth * 2
		sha.LP.Height = sha.ScreenHeight
		LW, LH, HLW, HLH := getLevelDimensionsInt()
		SetWorldImage(LW, LH)

		bg := com.NewBackground(background, 0, 0, 0, LW, LH, V)
		px, py = HLW-300, HLH+100
		player = com.NewPlayer(sha.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, sha.Red50)
		tester := com.NewCollideTest(sha.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, sha.Green50)
		wallT := com.NewWall(sha.IDWall, 0, 0, LW-size, size, sha.Blue50)
		wallL := com.NewWall(sha.IDWall, 0, size, size, LH-size, sha.Blue50)
		wallB := com.NewWall(sha.IDWall, size, LH-size, LW-size, size, sha.Blue50)
		wallR := com.NewWall(sha.IDWall, LW-size, 0, size, LH-size, sha.Blue50)
		wallCL := com.NewWall(sha.IDWall, 300, HLH-size/2, HLW-300, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-300, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 300-size/2, size, HLH-300, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-300, sha.Purple50)
		// add to interface lists
		DrawWorldList = append(DrawWorldList, &bg, &tb, &player, &tester, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &player, &tester, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &tester)
		CollideList = append(CollideList, &player, &tester)
		spwanRandomSquares(HitAbleList, 8, 50)

	} else if name == "lvl03" {
		sha.LP.Gravity = 0.04
		sha.LP.Width = sha.ScreenWidth
		sha.LP.Height = sha.ScreenHeight
		LW, LH, HLW, HLH := getLevelDimensionsInt()
		SetWorldImage(LW, LH)

		bg := com.NewBackground(background, 0, 0, 0, LW, LH, V)
		px, py = HLW, HLH+100
		player = com.NewPlayer(sha.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, sha.Red50)
		tester := com.NewCollideTest(sha.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, sha.Green50)
		wallB := com.NewWall(sha.IDWall, 0, LH-size, LW, size, sha.Green50)
		wallCL := com.NewWall(sha.IDWall, 400, HLH-size/2, HLW-400, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-400, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 400-size/2, size, HLH-400, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-400, sha.Purple50)
		// add to interface lists
		DrawWorldList = append(DrawWorldList, &bg, &tb, &player, &tester, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &player, &tester, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &tester)
		CollideList = append(CollideList, &player, &tester)
	}

	// update player position
	sha.LP.PlayerStartX = px
	sha.LP.PlayerStartY = py
}

func preloadImages() {
	var err error
	background, _, err = ebitenutil.NewImageFromFile("assets/space.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	spaceShip, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

// spwans squares in playable level dimensions, and makes sure that the squares dont overlap other objects
func spwanRandomSquares(list []com.HitAble, count, size int) {
	for i := 0; i < count; i++ {
		// get random free position and random velocity
		x, y := getRandonPosition(size, size, size, list)
		choices := []int{2, -2}
		vx := float64(choices[rand.Int()%len(choices)])
		vy := float64(choices[rand.Int()%len(choices)])
		s := com.NewSquare(sha.IDSquare, x, y, 0, com.NewVector(vx, vy), 0, 0, size, size, sha.Purple50)
		// add Square to interface lists
		DrawWorldList = append(DrawWorldList, &s)
		UpdateList = append(UpdateList, &s)
		CollideList = append(CollideList, &s)
		HitAbleList = append(HitAbleList, &s)
		//list = append(list, &s)
	}
}

func getLevelDimensionsInt() (int, int, int, int) {
	LW := int(sha.LP.Width)
	LH := int(sha.LP.Height)
	return LW, LH, int(LW / 2), int(LH / 2)
}

func getRandonPosition(offsetX, offsetY, space int, dontOverlap []com.HitAble) (int, int) {
	x := rand.Intn(int(sha.LP.Width)-(offsetX*2)) + offsetX
	y := rand.Intn(int(sha.LP.Height)-(offsetY*2)) + offsetY
	r := com.NewRect(x, y, space, space)
	for _, o := range dontOverlap {
		if o.GetObject().GetSolid() {
			if com.CheckOverlap(&r, o.GetObject().GetRect()) {
				x, y = getRandonPosition(offsetX, offsetX, space, dontOverlap)
			}
		}
	}
	return x, y
}
