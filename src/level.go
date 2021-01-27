package src

import (
	//"image/color"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	ass "moonlander/assets"
	com "moonlander/src/components"
	sha "moonlander/src/shared"

	"github.com/lafriks/go-tiled"
)

// ObjectTypes ..
type ObjectTypes struct {
	XMLName    xml.Name     `xml:"objecttypes"`
	ObjectType []ObjectType `xml:"objecttype"`
}

// ObjectType ..
type ObjectType struct {
	XMLName    xml.Name   `xml:"objecttype"`
	Name       string     `xml:"name,attr"`
	Color      string     `xml:"color,attr"`
	Properties []Property `xml:"property"`
}

// Property ..
type Property struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

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

// ClearLevel global variables
func ClearLevel() {
	DrawWorldList = nil
	DrawScreenList = nil
	HitAbleList = nil
	UpdateList = nil
	CollideList = nil
}

// LoadLevel loads a specific level
func LoadLevel(name string) {
	// xml created by Tiled with default values of object types
	objectTypePath := "assets/tiled/objecttypes.xml"

	if name == "lvl01" {
		loadTMX("assets/tiled/level01.tmx", objectTypePath)
		finalizeLevel()
	} else if name == "lvl02" {
		loadTMX("assets/tiled/level02.tmx", objectTypePath)
		finalizeLevel()
		spwanRandomSquares(HitAbleList, 8, 50)
	} else if name == "lvl03" {
		loadTMX("assets/tiled/level03.tmx", objectTypePath)
		finalizeLevel()
	}

	// anim example
	anim := com.NewAnimFromByte(ass.Runner_png, 200, 200, 0, com.Vector{}, com.NewFrame(0, 32, 32, 32, 8, 5))
	DrawWorldList = append(DrawWorldList, &anim)
	UpdateList = append(UpdateList, &anim)
}

func loadTMX(mapPath string, objectpath string) {
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
		BG:       getLevelBackground(m),
	}
	bg := com.NewBackground(sha.IDBG, sha.LP.BG, 0, 0, 0, sha.LP.Width, sha.LP.Height, com.Vector{})
	DrawWorldList = append(DrawWorldList, &bg)
	printLevelProperties()

	// Get object types properties default values
	objectTypes := getObjectTypes(objectpath)

	//loop through tile layers (not used atm)
	for _, layer := range m.Layers {
		for i, tile := range layer.Tiles {
			x, y := layer.GetTilePosition(i)
			if tile.Tileset != nil {
				id := strconv.FormatUint(uint64(tile.ID), 10)
				addLevelItem(id, "", x, y, m.TileWidth, m.TileHeight, 0,
					tile.Tileset.Tiles[tile.ID].Properties, objectTypes)
			}
		}
	}
	// loop through object layers
	for _, objLayer := range m.ObjectGroups {
		for _, obj := range objLayer.Objects {
			addLevelItem(obj.Type, obj.Name, int(obj.X), int(obj.Y), int(obj.Width),
				int(obj.Height), int(obj.Rotation), obj.Properties, objectTypes)
		}
	}
}

// Object properties are not in TMX level data
// unless an object instance overrides its default value or adds an instance properyu
func getObjectTypes(xmlPath string) []ObjectType {
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	// read xmlFile as a byte array and unmarshall to struct
	data, _ := ioutil.ReadAll(xmlFile)
	var objectTypes ObjectTypes
	xml.Unmarshal(data, &objectTypes)
	return objectTypes.ObjectType
}

func getObjectProps(id string, objectTypes []ObjectType) []Property {
	// match id, with objectType name
	// if it matches return Properties
	for _, o := range objectTypes {
		if o.Name == id {
			return o.Properties
		}
	}
	// TODO also add overridden values from TMX and add unique instance properties
	return nil
}

func addLevelItem(id, name string, x, y, w, h, rotation int, prop tiled.Properties, objectTypes []ObjectType) {
	//fmt.Printf("id:%v, name:%v, x:%v ,y:%v, w:%v, h:%v, r:%v, prop:%v", id, name, x, y, w, h, rotation, prop)

	p := getObjectProps(id, objectTypes)
	fmt.Println(p)

	// TODO want to use the objectTypes to determine in which intrfaces list's it should be appended
	// issue is that we can cast specific instances like a Player to something generic
	// and then add it to interface list, which it fullfills the contract
	// solutions
	// - interface stubs?
	// - type assertion?

	if id == "wall" {
		o := com.NewWall(sha.IDWall, x, y, w, h, sha.Blue50)
		DrawWorldList = append(DrawWorldList, &o)
		HitAbleList = append(HitAbleList, &o)
	}

	if id == "player" {
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
	tb := com.NewTextBlock(sha.ScreenWidth-80, 50)
	DrawScreenList = append(DrawScreenList, &tb)

	// add all checkpoints to finish
	finish.Checkpoints = checkpoints
	printLevelObjects()
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

func printLevelProperties() {
	fmt.Printf("level properties:\n%+v\n", sha.LP)
}

func printLevelObjects() {
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
