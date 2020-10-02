package constants

import (
	"image/color"
	"math"
)

// Some default colors (cant be constants, mutable, but we use them as constants)
var (
	Red    = color.RGBA{255, 0, 0, 128}
	Green  = color.RGBA{0, 255, 0, 128}
	Blue   = color.RGBA{0, 0, 255, 128}
	Purple = color.RGBA{255, 0, 255, 128}
	Yellow = color.RGBA{255, 255, 0, 128}

	// translate ids to name string
	ID = map[int]string{
		1: "player",
		2: "square",
		3: "wall",
		4: "tester",
	}
)

// Constants
const (
	ScreenWidth  = 1280
	ScreenHeight = 1024

	RadToDeg = 180 / math.Pi
	DegToRad = math.Pi / 180

	IDPlayer = 1
	IDSquare = 2
	IDWall   = 3
	IDTester = 4
)
