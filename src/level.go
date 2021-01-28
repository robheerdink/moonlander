package src

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

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
	DrawWorldList  []com.GameObject
	DrawScreenList []com.GameObject
	HitAbleList    []com.GameObject
	UpdateList     []com.GameObject
	CollideList    []com.GameObject
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
	checkpoints = nil
}

// LoadLevel loads a specific level
func LoadLevel(name string) {
	// xml created by Tiled with default values of object types
	// as long as the default values are not overriden, they will not be in TMX file
	objectTypePath := "assets/tiled/objecttypes.xml"
	if name == "lvl01" {
		loadTiledData("assets/tiled/level01.tmx", objectTypePath)
		finalizeLevel()
	} else if name == "lvl02" {
		loadTiledData("assets/tiled/level02.tmx", objectTypePath)
		finalizeLevel()
		spwanRandomSquares(HitAbleList, 8, 50)
	} else if name == "lvl03" {
		loadTiledData("assets/tiled/level03.tmx", objectTypePath)
		finalizeLevel()
	}
}

// Handles Tiled data
func loadTiledData(mapPath string, objectpath string) {
	m, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	// set Level properties from tmx map properties
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
	fmt.Printf("\n\nLevel: %v\nProperties:%+v\n\n", mapPath, sha.LP)

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

// Tiled object types xml to Structs
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

// Factory for populating the level with GameObjects
func addLevelItem(itemType, name string, x, y, w, h, rotation int, props tiled.Properties, objectTypes []ObjectType) {
	// fmt.Printf("id:%v, name:%v, x:%v ,y:%v, w:%v, h:%v, r:%v, prop:%v", id, name, x, y, w, h, rotation, prop)
	// get item properties
	p := getItemProps(itemType, props, objectTypes)
	switch itemType {
	case "wall":
		o := com.NewWall(sha.IDWall, x, y, w, h, sha.Blue50)
		addItemToList(&o, p)
		break
	case "player":
		player = com.NewPlayer(sha.IDPlayer, x, y, 0, com.Vector{}, 8, 8, 30, 48, sha.Red50)
		addItemToList(&player, p)
		break
	case "tester":
		o := com.NewCollideTest(sha.IDTester, x, y, 0, com.Vector{}, 4, 4, 24, 56, sha.Green50)
		addItemToList(&o, p)
		break
	case "cp":
		o := com.NewCheckpoint(sha.IDCheckpoint, x, y, w, h, sha.Cyan25, true)
		checkpoints = append(checkpoints, &o)
		addItemToList(&o, p)
		break
	case "finish":
		finish = com.NewFinish(sha.IDFinish, x, y, w, h, sha.White25, nil)
		addItemToList(&finish, p)
		break
	}
}

// Get the default, unique and overridden properties
func getItemProps(typ string, props tiled.Properties, objectTypes []ObjectType) []Property {
	// match id, with objectType.Name, if it matches return Properties
	for _, o := range objectTypes {
		if o.Name == typ {
			return o.Properties
		}
	}
	// TODO should check if values in props override the default value
	return nil
}

// Add the GameObjects to the correct lists, based on the (default) properties in Tiled
func addItemToList(item com.GameObject, properties []Property) {
	for _, p := range properties {
		if p.Name == "draw" && p.Default == "1" {
			DrawWorldList = append(DrawWorldList, item)
		}
		if p.Name == "hit" && p.Default == "1" {
			HitAbleList = append(HitAbleList, item)
		}
		if p.Name == "update" && p.Default == "1" {
			UpdateList = append(UpdateList, item)
		}
		if p.Name == "collide" && p.Default == "1" {
			CollideList = append(CollideList, item)
		}
	}
}

// Finilize level, do stuff we can only do when we have all objects or data
func finalizeLevel() {
	// set world render image (same size as the level)
	SetWorldImage(sha.LP.Width, sha.LP.Height)

	// player init position
	sha.LP.PlayerStartX = int(player.X)
	sha.LP.PlayerStartY = int(player.Y)

	// create gui
	tb := com.NewTextBlock(10, 24)
	DrawScreenList = append(DrawScreenList, &tb)

	// add all checkpoints to finish
	finish.Checkpoints = checkpoints
	printLevelObjects()
}

// Spawns squares in playable level dimensions, and makes sure that the squares dont overlap other objects
func spwanRandomSquares(list []com.GameObject, count, size int) {
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
	return m.Properties.GetFloat("gravity")
}

func getLevelFriction(m *tiled.Map) float64 {
	return m.Properties.GetFloat("friction")
}

func getLevelMaxMaps(m *tiled.Map) int {
	return m.Properties.GetInt("maxLaps")
}

func getRandonPosition(offsetX, offsetY, space int, dontOverlap []com.GameObject) (int, int) {
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
