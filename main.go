package main

import (
	con "moonlander/constants"
	gui "moonlander/gui"
	lvl "moonlander/level"

	"github.com/hajimehoshi/ebiten"
)

// Mode values (0,1,2)
const (
	ModeTitle int = iota
	ModeGame
	ModeGameOver
)

// Game implements ebiten.Game interface.
type Game struct {
	mode int
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
		for _, i := range lvl.UpdateList {
			i.Update(screen)
		}
		for _, i := range lvl.CollideList {
			i.Collide(lvl.ObjectList)
		}

		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			g.mode = ModeGameOver
			loadState(g, "")
		}

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
		for _, i := range lvl.DrawList {
			i.Draw(screen)
		}
	case ModeGameOver:
		gui.DrawGameOver(screen)
	}
}

// Load different modes of the game
func loadState(g *Game, action string) {
	if g.mode == ModeTitle {
		lvl.ClearLevel()
		gui.InitTitle()

	} else if g.mode == ModeGame {
		gui.ClearTitle()
		lvl.LoadLevel(action)

	} else if g.mode == ModeGameOver {
		lvl.ClearLevel()
		gui.InitGameOver()
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return con.ScreenWidth, con.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(con.ScreenWidth, con.ScreenHeight)
	ebiten.SetWindowTitle("Moon Lander!!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
