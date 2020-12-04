package src

import (
	//"image/color"
	"bytes"
	"image"
	"log"
	"math/rand"
	"time"

	ass "moonlander/assets"
	com "moonlander/src/components"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Variables related to level
var (
	background  *ebiten.Image
	spaceShip   *ebiten.Image
	player      com.Player
	DrawList    []com.Drawer
	HitAbleList []com.HitAble
	UpdateList  []com.Updater
	CollideList []com.Collider
)

// ClearLevel clears are references from level
func ClearLevel() {
	background = nil
	spaceShip = nil
	DrawList = nil
	HitAbleList = nil
	UpdateList = nil
	CollideList = nil
}

// LoadLevel loads a specific level
func LoadLevel(name string) {

	// common for all levels
	preloadImages()

	sha.WP = sha.WorldProperties{
		Gravity:     0.000,
		Friction:    0.992,
		LevelWidth:  1150,
		LevelHeight: 864,
	}
	sha.LP = sha.LevelProperties{
		Lives:      3,
		CurrentLap: 0,
		MaxLaps:    3,
		LapTimes:   []time.Duration{},
	}

	// some abbrivations
	LW := int(sha.WP.LevelWidth)
	LH := int(sha.WP.LevelHeight)
	HLW := int(LW / 2)
	HLH := int(LH / 2)
	V := com.Vector{}
	var size int = 16
	var cpSize int = 4

	bg := com.NewSprite(background, 0, 0, 0, V)
	tb := com.NewTextBlock(LW-150, 50)
	px, py := 0, 0

	img, _, err := image.Decode(bytes.NewReader(ass.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	if name == "lvl01" {
		runnerImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		anim := com.NewAnim(runnerImage, 200, 200, 0, V, com.NewFrame(0, 32, 32, 32, 8, 5))

		sha.WP.Gravity = 0.015
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
		wallCL := com.NewWall(sha.IDWall, 200, HLH-size/2, HLW-200, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-200, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 200-size/2, size, HLH-200, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-200, sha.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &anim, &tb, &player, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &player, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &anim)
		CollideList = append(CollideList, &player)
		//spwanRandomSquares(HitAbleList, 2, 50)

	} else if name == "lvl02" {
		sha.WP.Gravity = 0.00
		px, py = HLW-300, HLH+100
		player := com.NewPlayer(sha.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, sha.Red50)
		tester := com.NewCollideTest(sha.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, sha.Green50)
		wallT := com.NewWall(sha.IDWall, 0, 0, LW-size, size, sha.Blue50)
		wallL := com.NewWall(sha.IDWall, 0, size, size, LH-size, sha.Blue50)
		wallB := com.NewWall(sha.IDWall, size, LH-size, LW-size, size, sha.Blue50)
		wallR := com.NewWall(sha.IDWall, LW-size, 0, size, LH-size, sha.Blue50)
		wallCL := com.NewWall(sha.IDWall, 200, HLH-size/2, HLW-200, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-200, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 200-size/2, size, HLH-200, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-200, sha.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &tb, &player, &tester, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &player, &tester, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &tester)
		CollideList = append(CollideList, &player, &tester)
		spwanRandomSquares(HitAbleList, 500, 8)

	} else if name == "lvl03" {
		sha.WP.Gravity = 0.04
		px, py = HLW, HLH+100
		player = com.NewPlayer(sha.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, sha.Red50)
		tester := com.NewCollideTest(sha.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, sha.Green50)
		wallB := com.NewWall(sha.IDWall, 0, LH-size, LW, size, sha.Green50)
		wallCL := com.NewWall(sha.IDWall, 400, HLH-size/2, HLW-400, size, sha.Cyan50)
		wallCR := com.NewWall(sha.IDWall, HLW, HLH-size/2, HLW-400, size, sha.Yellow50)
		wallCT := com.NewWall(sha.IDWall, HLW-size/2, 400-size/2, size, HLH-400, sha.Green50)
		wallCB := com.NewWall(sha.IDWall, HLW-size/2, HLH+size/2, size, HLH-400, sha.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &tb, &player, &tester, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &player, &tester, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &player, &tester)
		CollideList = append(CollideList, &player, &tester)
	}

	// update player position
	sha.LP.PX = px
	sha.LP.PY = py
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
		DrawList = append(DrawList, &s)
		UpdateList = append(UpdateList, &s)
		CollideList = append(CollideList, &s)
		HitAbleList = append(HitAbleList, &s)
		//list = append(list, &s)
	}
}

func getRandonPosition(offsetX, offsetY, space int, dontOverlap []com.HitAble) (int, int) {
	x := rand.Intn(int(sha.WP.LevelWidth)-(offsetX*2)) + offsetX
	y := rand.Intn(int(sha.WP.LevelHeight)-(offsetY*2)) + offsetY
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
