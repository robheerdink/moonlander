package level

import (
	"log"
	"math/rand"

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
	UpdateList  []comp.Updater
	CollideList []comp.Collider
	ObjectList  []*comp.Object
)

// ClearLevel clears are references from level
func ClearLevel() {
	background = nil
	spaceShip = nil
	DrawList = nil
	UpdateList = nil
	CollideList = nil
	ObjectList = nil
}

// LoadLevel loads a specific level
func LoadLevel(name string) {

	// common for all levels
	preloadImages()

	comp.WP = comp.WorldProperties{
		Gravity:     0,
		Friction:    0.992,
		LevelWidth:  1280,
		LevelHeight: 1024,
	}
	HLW := float64(comp.WP.LevelWidth / 2)
	HLH := float64(comp.WP.LevelHeight / 2)
	V := comp.NewVector(0, 0)
	var walls []*comp.Object
	var size int = 16

	if name == "lvl01" {
		// create components
		bg := comp.NewSprite(0, background, 0, 0, 0, V)
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, HLW-300, HLH+100, 0, V, 18, 18, 28, 28, con.Red)
		p2 := comp.NewCollideTest(con.IDTester, HLW-200, HLH+200, 0, V, 0, 0, 16, 16, con.Green)
		wallT := comp.NewObjectNoImage(con.IDWall, 0, 0, 0, V, 0, 0, comp.WP.LevelWidth-size, size, con.Blue)
		wallL := comp.NewObjectNoImage(con.IDWall, 0, float64(size), 0, V, 0, 0, size, comp.WP.LevelHeight-size, con.Blue)
		wallB := comp.NewObjectNoImage(con.IDWall, 16, float64(comp.WP.LevelHeight)-float64(size), 0, V, 0, 0, comp.WP.LevelWidth, size, con.Blue)
		wallR := comp.NewObjectNoImage(con.IDWall, float64(comp.WP.LevelWidth)-float64(size), 0, 0, V, 0, 0, size, comp.WP.LevelHeight-size, con.Blue)
		wallCL := comp.NewObjectNoImage(con.IDWall, 400, HLH-float64(size), 0, V, 0, 0, int(HLW-400), size, con.Blue)
		wallCR := comp.NewObjectNoImage(con.IDWall, float64(HLW), HLH-float64(size), 0, V, 0, 0, int(HLW-400), size, con.Yellow)
		wallCT := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(400-size), 0, V, 0, 0, size, int(HLH-400), con.Green)
		wallCB := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(HLH), 0, V, 0, 0, size, int(HLH)-400, con.Purple)
		walls = append(walls, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)

		// add to interface lists
		DrawList = append(DrawList, &bg, &p1, &p2)
		UpdateList = append(UpdateList, &p1, &p2)
		CollideList = append(CollideList, &p1, &p2)
		for _, o := range walls {
			DrawList = append(DrawList, o)
		}

		var dontOverlap []*comp.Object = append(walls, &p1.Object, &p2.Object)
		spwanRandomSquares(dontOverlap, 5, 100)

		// add all object to which we can collide to Objectlist
		for _, o := range CollideList {
			ObjectList = append(ObjectList, o.GetObject())
		}
		for _, o := range walls {
			ObjectList = append(ObjectList, o)
		}

	} else if name == "lvl02" {

		// create components
		bg := comp.NewSprite(0, background, 0, 0, 0, V)
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, HLW+100, HLH+100, 0, V, 18, 18, 28, 28, con.Red)
		p2 := comp.NewCollideTest(con.IDTester, HLW-100, HLH+100, 0, V, 4, 4, 24, 24, con.Green)
		wallT := comp.NewObjectNoImage(con.IDWall, 0, 0, 0, V, 0, 0, comp.WP.LevelWidth-size, size, con.Blue)
		wallL := comp.NewObjectNoImage(con.IDWall, 0, float64(size), 0, V, 0, 0, size, comp.WP.LevelHeight-size, con.Blue)
		wallB := comp.NewObjectNoImage(con.IDWall, 16, float64(comp.WP.LevelHeight)-float64(size), 0, V, 0, 0, comp.WP.LevelWidth, size, con.Blue)
		wallR := comp.NewObjectNoImage(con.IDWall, float64(comp.WP.LevelWidth)-float64(size), 0, 0, V, 0, 0, size, comp.WP.LevelHeight-size, con.Blue)
		wallCL := comp.NewObjectNoImage(con.IDWall, 300, HLH-float64(size), 0, V, 0, 0, int(HLW-300), size, con.Blue)
		wallCR := comp.NewObjectNoImage(con.IDWall, float64(HLW), HLH-float64(size), 0, V, 0, 0, int(HLW-300), size, con.Yellow)
		wallCT := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(300-size), 0, V, 0, 0, size, int(HLH-300), con.Green)
		wallCB := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(HLH), 0, V, 0, 0, size, int(HLH)-300, con.Purple)
		walls = append(walls, &wallT, &wallL, &wallB, &wallR, &wallCL, &wallCR, &wallCT, &wallCB)

		// add to interface lists
		DrawList = append(DrawList, &bg, &p1, &p2)
		UpdateList = append(UpdateList, &p1, &p2)
		CollideList = append(CollideList, &p1, &p2)
		for _, o := range walls {
			DrawList = append(DrawList, o)
		}

		var dontOverlap []*comp.Object = append(walls, &p1.Object, &p2.Object)
		spwanRandomSquares(dontOverlap, 500, 8)

		// add all object to which we can collide to Objectlist
		for _, o := range CollideList {
			ObjectList = append(ObjectList, o.GetObject())
		}
		for _, o := range walls {
			ObjectList = append(ObjectList, o)
		}

	} else if name == "lvl03" {
		// set world properties
		comp.WP.Gravity = 0.04

		// create components
		bg := comp.NewSprite(0, background, 0, 0, 0, V)
		p1 := comp.NewPlayer(con.IDPlayer, spaceShip, HLW+100, HLH+100, 0, V, 18, 18, 28, 28, con.Red)
		wallB := comp.NewObjectNoImage(con.IDWall, 16, float64(comp.WP.LevelHeight)-float64(size), 0, V, 0, 0, comp.WP.LevelWidth, size, con.Blue)
		wallCL := comp.NewObjectNoImage(con.IDWall, 300, HLH-float64(size), 0, V, 0, 0, int(HLW-300), size, con.Blue)
		wallCR := comp.NewObjectNoImage(con.IDWall, float64(HLW), HLH-float64(size), 0, V, 0, 0, int(HLW-300), size, con.Yellow)
		wallCT := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(300-size), 0, V, 0, 0, size, int(HLH-300), con.Green)
		wallCB := comp.NewObjectNoImage(con.IDWall, HLW-float64(size), float64(HLH), 0, V, 0, 0, size, int(HLH)-300, con.Purple)
		walls = append(walls, &wallB, &wallCL, &wallCR, &wallCT, &wallCB)

		// add to interface lists
		DrawList = append(DrawList, &bg, &p1)
		UpdateList = append(UpdateList, &p1)
		CollideList = append(CollideList, &p1)
		for _, o := range walls {
			DrawList = append(DrawList, o)
		}

		// add all object to which we can collide to Objectlist
		for _, o := range CollideList {
			ObjectList = append(ObjectList, o.GetObject())
		}
		for _, o := range walls {
			ObjectList = append(ObjectList, o)
		}
	}
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
func spwanRandomSquares(list []*comp.Object, count, size int) {
	for i := 0; i < count; i++ {
		// get random free position and random velocity
		x, y := getRandonPosition(size, size, size, list)
		choices := []int{2, -2}
		vx := float64(choices[rand.Int()%len(choices)])
		vy := float64(choices[rand.Int()%len(choices)])
		s := comp.NewSquare(con.IDSquare, x, y, 0, comp.NewVector(vx, vy), 0, 0, size, size, con.Purple)
		list = append(list, &s.Object)
		DrawList = append(DrawList, &s)
		UpdateList = append(UpdateList, &s)
		CollideList = append(CollideList, &s)
	}
}

func getRandonPosition(offsetX, offsetY, space int, dontOverlap []*comp.Object) (float64, float64) {
	x := rand.Intn(int(comp.WP.LevelWidth)-(offsetX*2)) + offsetX
	y := rand.Intn(int(comp.WP.LevelHeight)-(offsetY*2)) + offsetY
	r := comp.Rect{X: x, Y: y, W: space, H: space}
	xf, yf := float64(x), float64(y)
	for _, o := range dontOverlap {
		if comp.CheckOverlap(&r, &o.Rect) {
			xf, yf = getRandonPosition(offsetX, offsetX, space, dontOverlap)
		}
	}
	return xf, yf
}
