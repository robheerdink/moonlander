package src

import (
	"fmt"
	gui "moonlander/src/gui"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

// Mode values (0,1,2)
const (
	ModeTitle int = iota
	ModeGame
	ModeGameOver
)

// Game implements ebiten.Game interface.
type Game struct {
	mode   int
	world  *ebiten.Image
	camera Camera
}

// Run this code once at startup app
func init() {
	gui.InitTitle()
}

// Update proceeds the game state.
func (g *Game) Update(screen *ebiten.Image) error {
	var action string

	switch g.mode {
	case ModeTitle:
		action = gui.UpdateTitle(screen)
		if action != "" {
			g.mode = ModeGame
			loadState(g, action)
		}

	case ModeGame:
		// loop through update list and collide list

		for _, i := range UpdateList {
			i.Update(screen)
		}
		for _, i := range CollideList {
			i.Collide(HitAbleList)
		}

		// test game-over screen
		if ebiten.IsKeyPressed(ebiten.KeyBackslash) {
			g.mode = ModeGameOver
			loadState(g, "")
		}

		// camera test controls
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.camera.Position[0] -= 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.camera.Position[0] += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.camera.Position[1] -= 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.camera.Position[1] += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			if g.camera.ZoomFactor > -240 {
				g.camera.ZoomFactor -= 2
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			if g.camera.ZoomFactor < 240 {
				g.camera.ZoomFactor += 2
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.camera.Rotation += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.camera.Reset()
		}

		//- default camera postion is 0,0
		//- how to get player x, y, this seems strange
		//- some paint clear issue
		//- look in to draw world image size, you can e.g fly vissble out of level on right side 2000x2000 world image, not on left
		//- need to seperate draw in world and draw on screen for gui
		//-
		x, y, _, _ := player.GetImageInfo()
		fmt.Println(x, y)
		fmt.Println(g.camera)
		g.camera.Position[0] = x - sha.ScreenWidth/2
		g.camera.Position[1] = y - sha.ScreenHeight/2

	case ModeGameOver:
		action = gui.UpdateGameOver(screen)
		if action != "" {
			g.mode = ModeTitle
			loadState(g, action)
		}
	}

	// handle escape in game or gameover screen
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if g.mode == ModeGame || g.mode == ModeGameOver {
			g.mode = ModeTitle
			loadState(g, "")
		}
	}
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case ModeTitle:
		gui.DrawTitle(screen)
	case ModeGame:

		// draw in world
		for _, i := range DrawList {
			i.Draw(g.world)
		}

		// draw world to screen
		// op := &ebiten.DrawImageOptions{}
		// screen.DrawImage(g.world, op)

		// render world in camera
		g.camera.Render(g.world, screen)

	case ModeGameOver:
		gui.DrawGameOver(screen)
	}
}

// Load different modes of the game
func loadState(g *Game, action string) {
	if g.mode == ModeTitle {
		ClearLevel()
		gui.InitTitle()

	} else if g.mode == ModeGame {
		gui.ClearTitle()
		LoadLevel(action)

	} else if g.mode == ModeGameOver {
		ClearLevel()
		gui.InitGameOver()
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return sha.ScreenWidth, sha.ScreenHeight
}

func Start() {
	ebiten.SetWindowSize(sha.ScreenWidth, sha.ScreenHeight)
	ebiten.SetWindowTitle("Moon Lander!!")

	// set camera and world
	g := &Game{}
	g.camera = Camera{ViewPort: f64.Vec2{sha.ScreenWidth, sha.ScreenHeight}}
	g.world, _ = ebiten.NewImage(2000, 2000, ebiten.FilterDefault)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
