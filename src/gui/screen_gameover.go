package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	gameOverText = `Gameover, press left mouse or escape.`
)

// ClearGameOver clears gameover screen
func ClearGameOver() {
}

// InitGameOver inits gameover screen
func InitGameOver() {
}

// UpdateGameOver updates gameover screen
func UpdateGameOver(screen *ebiten.Image) string {
	mouse.update()
	if mouse.pressed {
		return "title"
	}
	return ""
}

// DrawGameOver draws gameover screen
func DrawGameOver(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 64, 64, 255})
	text.Draw(screen, gameOverText, fontBig, 20, 160, color.White)
}
