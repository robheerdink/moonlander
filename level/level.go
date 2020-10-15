package level

import (
	//"image/color"
	"bytes"
	"image"
	"log"
	"math/rand"
	"time"

	images "moonlander/assets"
	comp "moonlander/components"
	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Variables related to level
var (
	background  *ebiten.Image
	spaceShip   *ebiten.Image
	DrawList    []comp.Drawer
	HitAbleList []comp.HitAble
	UpdateList  []comp.Updater
	CollideList []comp.Collider
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

var (
	runnerImage *ebiten.Image
)

// LoadLevel loads a specific level
func LoadLevel(name string) {

	// common for all levels
	preloadImages()

	comp.WP = comp.WorldProperties{
		Gravity:     0.000,
		Friction:    0.992,
		LevelWidth:  1150,
		LevelHeight: 864,
	}
	comp.LP = comp.LevelProperties{
		Lives:      3,
		CurrentLap: 0,
		MaxLaps:    3,
		LapTimes:   []time.Duration{},
	}

	// some abbrivations
	LW := int(comp.WP.LevelWidth)
	LH := int(comp.WP.LevelHeight)
	HLW := int(LW / 2)
	HLH := int(LH / 2)
	V := comp.Vector{}
	var size int = 16
	var cpSize int = 4

	bg := comp.NewSprite(background, 0, 0, 0, V)
	tb := comp.NewTextBlock(LW-150, 50)
	px, py := 0, 0

	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	if name == "lvl01" {
		runnerImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		anim := comp.NewAnim(runnerImage, 200, 200, 0, V, comp.NewFrame(0, 32, 32, 32, 8, 5))

		comp.WP.Gravity = 0 //0.015
		px, py = HLW-300, HLH+100
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, con.Red50)
		cp1 := comp.NewCheckpoint(con.IDCheckpoint, HLW-cpSize, HLH-cpSize/2, HLW, cpSize, con.Cyan25, true)
		cp2 := comp.NewCheckpoint(con.IDCheckpoint, HLW-cpSize/2, 0, cpSize, HLH, con.Green25, true)
		cp3 := comp.NewCheckpoint(con.IDCheckpoint, 0, HLH-cpSize/2, HLW, cpSize, con.Yellow25, true)
		finish := comp.NewFinish(con.IDFinish, HLW-cpSize/2, HLH, cpSize, HLH, con.White25, []*comp.Checkpoint{&cp1, &cp2, &cp3})
		wallT := comp.NewWall(con.IDWall, 0, 0, LW-size, size, con.Blue50)
		wallL := comp.NewWall(con.IDWall, 0, size, size, LH-size, con.Blue50)
		wallB := comp.NewWall(con.IDWall, size, LH-size, LW-size, size, con.Blue50)
		wallR := comp.NewWall(con.IDWall, LW-size, 0, size, LH-size, con.Blue50)
		wallCL := comp.NewWall(con.IDWall, 200, HLH-size/2, HLW-200, size, con.Cyan50)
		wallCR := comp.NewWall(con.IDWall, HLW, HLH-size/2, HLW-200, size, con.Yellow50)
		wallCT := comp.NewWall(con.IDWall, HLW-size/2, 200-size/2, size, HLH-200, con.Green50)
		wallCB := comp.NewWall(con.IDWall, HLW-size/2, HLH+size/2, size, HLH-200, con.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &anim, &tb, &p1, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &p1, &finish, &cp1, &cp2, &cp3, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &p1, &anim)
		CollideList = append(CollideList, &p1)
		//spwanRandomSquares(HitAbleList, 2, 50)

	} else if name == "lvl02" {
		comp.WP.Gravity = 0.00
		px, py = HLW-300, HLH+100
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, con.Red50)
		p2 := comp.NewCollideTest(con.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, con.Green50)
		wallT := comp.NewWall(con.IDWall, 0, 0, LW-size, size, con.Blue50)
		wallL := comp.NewWall(con.IDWall, 0, size, size, LH-size, con.Blue50)
		wallB := comp.NewWall(con.IDWall, size, LH-size, LW-size, size, con.Blue50)
		wallR := comp.NewWall(con.IDWall, LW-size, 0, size, LH-size, con.Blue50)
		wallCL := comp.NewWall(con.IDWall, 200, HLH-size/2, HLW-200, size, con.Cyan50)
		wallCR := comp.NewWall(con.IDWall, HLW, HLH-size/2, HLW-200, size, con.Yellow50)
		wallCT := comp.NewWall(con.IDWall, HLW-size/2, 200-size/2, size, HLH-200, con.Green50)
		wallCB := comp.NewWall(con.IDWall, HLW-size/2, HLH+size/2, size, HLH-200, con.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &tb, &p1, &p2, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &p1, &p2, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &p1, &p2)
		CollideList = append(CollideList, &p1, &p2)
		spwanRandomSquares(HitAbleList, 500, 8)

	} else if name == "lvl03" {
		comp.WP.Gravity = 0.03
		px, py = HLW, HLH+100
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, px, py, 0, V, 8, 8, 30, 48, con.Red50)
		p2 := comp.NewCollideTest(con.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 56, con.Green50)
		wallB := comp.NewWall(con.IDWall, 0, LH-size, LW, size, con.Green50)
		wallCL := comp.NewWall(con.IDWall, 400, HLH-size/2, HLW-400, size, con.Cyan50)
		wallCR := comp.NewWall(con.IDWall, HLW, HLH-size/2, HLW-400, size, con.Yellow50)
		wallCT := comp.NewWall(con.IDWall, HLW-size/2, 400-size/2, size, HLH-400, con.Green50)
		wallCB := comp.NewWall(con.IDWall, HLW-size/2, HLH+size/2, size, HLH-400, con.Purple50)
		// add to interface lists
		DrawList = append(DrawList, &bg, &tb, &p1, &p2, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		HitAbleList = append(HitAbleList, &p1, &p2, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)
		UpdateList = append(UpdateList, &p1, &p2)
		CollideList = append(CollideList, &p1, &p2)
	}

	// update player position
	comp.LP.PX = px
	comp.LP.PY = py
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
func spwanRandomSquares(list []comp.HitAble, count, size int) {
	for i := 0; i < count; i++ {
		// get random free position and random velocity
		x, y := getRandonPosition(size, size, size, list)
		choices := []int{2, -2}
		vx := float64(choices[rand.Int()%len(choices)])
		vy := float64(choices[rand.Int()%len(choices)])
		s := comp.NewSquare(con.IDSquare, x, y, 0, comp.NewVector(vx, vy), 0, 0, size, size, con.Purple50)
		// add Square to interface lists
		DrawList = append(DrawList, &s)
		UpdateList = append(UpdateList, &s)
		CollideList = append(CollideList, &s)
		HitAbleList = append(HitAbleList, &s)
		//list = append(list, &s)
	}
}

func getRandonPosition(offsetX, offsetY, space int, dontOverlap []comp.HitAble) (int, int) {
	x := rand.Intn(int(comp.WP.LevelWidth)-(offsetX*2)) + offsetX
	y := rand.Intn(int(comp.WP.LevelHeight)-(offsetY*2)) + offsetY
	r := comp.NewRect(x, y, space, space)
	for _, o := range dontOverlap {
		if o.GetObject().GetSolid() {
			if comp.CheckOverlap(&r, o.GetObject().GetRect()) {
				x, y = getRandonPosition(offsetX, offsetX, space, dontOverlap)
			}
		}
	}
	return x, y
}
