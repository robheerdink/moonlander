package gui

import (
	"fmt"
	"image/color"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	sampleText = `Moonlander`
)

var (
	btnList    []clickable
	clickedBtn clickable
)

// ClearTitle clears the title screen
func ClearTitle() {
	btnList = nil
	clickedBtn = nil
}

// InitTitle inits the title screen
func InitTitle() {
	// create some buttons
	w, h := 250, 100
	x1, x2, x3, y := sha.ScreenWidth/4-w/2, sha.ScreenWidth/2-w/2, sha.ScreenWidth/4*3-w/2, sha.ScreenHeight/3-h/2
	btnColor := color.RGBA{0, 255, 0, 128}
	txtColor := color.RGBA{0, 0, 0, 128}
	btn1 := newButton("lvl01", "Level 1 Amazing!!", x1, y, w, h, fontNormal, btnColor, txtColor)
	btn2 := newButton("lvl02", "Level 2", x2, y, w, h, fontNormal, btnColor, txtColor)
	btn3 := newButton("lvl03", "Level 3", x3, y, w, h, fontNormal, btnColor, txtColor)
	btnList = append(btnList, &btn1, &btn2, &btn3)
}

// UpdateTitle ..
func UpdateTitle(screen *ebiten.Image) string {
	mouse.update()
	if mouse.pressed {
		clickedBtn = checkHits(mouse.x, mouse.y, btnList)
		if clickedBtn != nil {
			return clickedBtn.getName()
		}
	}
	return ""
}

// DrawTitle ..
func DrawTitle(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	// Draw info
	msgFPS := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	msgMouse := fmt.Sprintf("X: %d, Y: %d Pressed: %t", mouse.x, mouse.y, mouse.pressed)
	text.Draw(screen, msgFPS, fontNormal, 20, 40, color.White)
	text.Draw(screen, msgMouse, fontNormal, 20, 80, color.White)
	text.Draw(screen, sampleText, fontNormal, 20, 120, color.White)

	text.Draw(screen, ""+
		"main menu     = esc \n\n"+
		"Move up       = up\n"+
		"Move right    = right\n"+
		"Move left     = left\n"+
		"Move down     = down\n"+
		"rotate left   = z\n"+
		"rotate right  = x\n"+
		"reset player  = backspace\n\n"+
		"camera move   = WSAD\n"+
		"camera rotate = Q/E\n"+
		"camera zoom   = SHIFT/CONTROL\n"+
		"camera reset  = SPACE\n\n"+
		"test obj      = T,G,F,H,R,Y",
		fontArcade, 200, 450, color.White)

	// draw all buttons
	for _, btn := range btnList {
		btn.draw(screen)
	}

}
