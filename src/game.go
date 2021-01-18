package src

import (
	gui "moonlander/src/gui"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

// Game implements ebiten.Game interface.
type Game struct {
	mode   int
	world  *ebiten.Image
	camera Camera
}

// Mode values (0,1,2)
const (
	ModeTitle int = iota
	ModeGame
	ModeGameOver
)

var (
	g *Game
)

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

		// update camera
		g.camera.Update()

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
		for _, i := range DrawWorldList {
			i.Draw(g.world)
		}
		// render world in camera
		g.camera.Render(g.world, screen)

		// draw on screen (gui)
		for _, i := range DrawScreenList {
			i.Draw(screen)
		}
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
		g.camera.Reset()

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

// Start the game
func Start() {
	ebiten.SetWindowSize(sha.ScreenWidth, sha.ScreenHeight)
	ebiten.SetWindowTitle("Moon Lander!!")

	// set camera
	g = &Game{}
	g.camera = Camera{ViewPort: f64.Vec2{sha.ScreenWidth, sha.ScreenHeight}}

	// Rungame starts main loop
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// SetWorldImage sets the size of world canvas image (should be same as level size)
// should be called when changing level
func SetWorldImage(width, height int) {
	g.world, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
}
