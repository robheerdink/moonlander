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

	ass "moonlander/assets"
	com "moonlander/src/components"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"github.com/lafriks/go-tiled"
)

// Variables related to level
var (
	player         com.Player
	DrawWorldList  []com.Drawer
	DrawScreenList []com.Drawer
	HitAbleList    []com.HitAble
	UpdateList     []com.Updater
	CollideList    []com.Collider
	checkpoints    []*com.Checkpoint
	finish         com.Finish
)

func ClearLevel() {
	DrawWorldList = nil
	DrawScreenList = nil
	HitAbleList = nil
	UpdateList = nil
	CollideList = nil
}

func LoadTMX(mapPath string) {
	/* LoadTMX loads a level with format TMX (tiled) */
	// parse tmx format
	m, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	// set Level properties from map properties
	sha.LP = sha.LevelProperties{
		Gravity:  getLevelGravity(m),
		Friction: getLevelFriction(m),
		MaxLaps:  getLevelMaxMaps(m),
		Width:    m.Width * m.TileWidth,
		Height:   m.Height * m.TileHeight,
	}
	bgPath := getLevelBackground(m)
	bg := com.NewBackground(sha.IDBG, bgPath, 0, 0, 0, sha.LP.Width, sha.LP.Height, com.Vector{})
	DrawWorldList = append(DrawWorldList, &bg)
	fmt.Println("level properties:")
	fmt.Println("sha.LP.Width    :", sha.LP.Width)
	fmt.Println("sha.LP.Height   :", sha.LP.Height)
	fmt.Println("sha.LP.Gravity  :", sha.LP.Gravity)
	fmt.Println("sha.LP.Friction :", sha.LP.Friction)
	fmt.Println("sha.LP.MaxLaps  :", sha.LP.MaxLaps)
	fmt.Println("background path :", bgPath)
	fmt.Println("----------------------------------")

	//loop through tile layers
	for _, layer := range m.Layers {
		for i, tile := range layer.Tiles {
			x, y := layer.GetTilePosition(i)
			if tile.Tileset != nil {
				levelItem(strconv.FormatUint(uint64(tile.ID), 10), "", x, y, m.TileWidth, m.TileHeight, 0, tile.Tileset.Tiles[tile.ID].Properties)
			}
		}
	}
	// loop through object layers
	for _, objLayer := range m.ObjectGroups {
		for _, obj := range objLayer.Objects {
			levelItem(obj.Type, obj.Name, int(obj.X), int(obj.Y), int(obj.Width), int(obj.Height), int(obj.Rotation), obj.Properties)
		}
	}
}

func levelItem(id, name string, x, y, w, h, rotation int, prop tiled.Properties) {
	fmt.Printf("id:%v, name:%v, x:%v ,y:%v, w:%v, h:%v, r:%v, prop:%v", id, name, x, y, w, h, rotation, prop)

	if id == "wall" {
		fmt.Printf(" >> add wall\n")
		o := com.NewWall(sha.IDWall, x, y, w, h, sha.Blue50)
		DrawWorldList = append(DrawWorldList, &o)
		HitAbleList = append(HitAbleList, &o)
	}
	if id == "player" {
		fmt.Printf(" >> add player\n")
		player = com.NewPlayer(sha.IDPlayer, x, y, 0, com.Vector{}, 8, 8, 30, 48, sha.Red50)
		DrawWorldList = append(DrawWorldList, &player)
		HitAbleList = append(HitAbleList, &player)
		UpdateList = append(UpdateList, &player)
		CollideList = append(CollideList, &player)
		sha.LP.PlayerStartX = x
		sha.LP.PlayerStartY = y
	}
	if id == "tester" {
		o := com.NewCollideTest(sha.IDTester, x, y, 0, com.Vector{}, 4, 4, 24, 56, sha.Green50)
		DrawWorldList = append(DrawWorldList, &o)
		HitAbleList = append(HitAbleList, &o)
		UpdateList = append(UpdateList, &o)
		CollideList = append(CollideList, &o)
	}
	if id == "cp" {
		o := com.NewCheckpoint(sha.IDCheckpoint, x, y, w, h, sha.Cyan25, true)
		checkpoints = append(checkpoints, &o)
		DrawWorldList = append(DrawWorldList, &o)
		HitAbleList = append(HitAbleList, &o)
	}
	if id == "finish" {
		finish = com.NewFinish(sha.IDFinish, x, y, w, h, sha.White25, nil)
		DrawWorldList = append(DrawWorldList, &finish)
		HitAbleList = append(HitAbleList, &finish)
	}
}

func finalizeLevel() {
	// set world render image (same size as the level)
	SetWorldImage(sha.LP.Width, sha.LP.Height)

	// create gui
	tb := com.NewTextBlock(sha.LP.Width-150, 50)
	DrawScreenList = append(DrawScreenList, &tb)

	// add all checkpoints to finish
	finish.Checkpoints = checkpoints
	printAllLevelObjects()
}

// LoadLevel loads a specific level
func LoadLevel(name string) {
	if name == "lvl01" {
		LoadTMX("assets/tiled/level01.tmx")
		finalizeLevel()
	} else if name == "lvl02" {
		LoadTMX("assets/tiled/level02.tmx")
		finalizeLevel()
		spwanRandomSquares(HitAbleList, 8, 50)
	} else if name == "lvl03" {
		LoadTMX("assets/tiled/level03.tmx")
		finalizeLevel()
	}

	// anim example
	img, _, err := image.Decode(bytes.NewReader(ass.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	anim := com.NewAnim(runnerImage, 200, 200, 0, com.Vector{}, com.NewFrame(0, 32, 32, 32, 8, 5))
	DrawWorldList = append(DrawWorldList, &anim)
	UpdateList = append(UpdateList, &anim)
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

func getLevelIndex(x, y int, m *tiled.Map) int {
	return x + y*m.Width
}

func getLevelXY(index int, m *tiled.Map) (int, int) {
	x := index % m.Width
	y := int(index / m.Width)
	return x, y
}

func getLevelBackground(m *tiled.Map) string {
	return m.Properties.GetString("background")
}

func getLevelGravity(m *tiled.Map) float64 {
	value, err := strconv.ParseFloat(m.Properties.GetString("gravity"), 64)
	if err != nil {
		fmt.Println("Failed to load Gravity from Map ", value, err)
	}
	return value
}

func getLevelFriction(m *tiled.Map) float64 {
	value, err := strconv.ParseFloat(m.Properties.GetString("friction"), 64)
	if err != nil {
		fmt.Println("Failed to load Friction from Map ", value, err)
	}
	return value
}

func getLevelMaxMaps(m *tiled.Map) int {
	value, err := strconv.Atoi(m.Properties.GetString("maxLaps"))
	if err != nil {
		fmt.Println("Failed to load MaxLaps from Map ", value, err)
	}
	return value
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

func printAllLevelObjects() {
	fmt.Println("\nDrawWorldList")
	for _, o := range DrawWorldList {
		fmt.Println(o.GetInfo())
	}
	fmt.Println("\nDrawScreenList")
	for _, o := range DrawScreenList {
		fmt.Println(o.GetInfo())
	}
	fmt.Println("\nHitAbleList")
	for _, o := range HitAbleList {
		fmt.Println(o.GetInfo())
	}
	fmt.Println("\nUpdateList")
	for _, o := range UpdateList {
		fmt.Println(o.GetInfo())
	}
	fmt.Println("\nCollideList")
	for _, o := range CollideList {
		fmt.Println(o.GetInfo())
	}
}
