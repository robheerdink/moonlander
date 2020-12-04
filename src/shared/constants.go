package shared

import (
	"image/color"
)

// Some default colors (cant be constants, mutable, but we use them as constants)
var (
	Red    = color.RGBA{255, 0, 0, 255}
	Green  = color.RGBA{0, 255, 0, 255}
	Blue   = color.RGBA{0, 0, 255, 255}
	Yellow = color.RGBA{255, 255, 0, 255}
	Cyan   = color.RGBA{0, 255, 255, 255}
	Purple = color.RGBA{255, 0, 255, 255}
	White  = color.RGBA{255, 255, 255, 255}

	Red50    = color.RGBA{255, 0, 0, 128}
	Green50  = color.RGBA{0, 255, 0, 128}
	Blue50   = color.RGBA{0, 0, 255, 128}
	Yellow50 = color.RGBA{255, 255, 0, 128}
	Cyan50   = color.RGBA{0, 255, 255, 128}
	Purple50 = color.RGBA{255, 0, 255, 128}
	White50  = color.RGBA{255, 255, 255, 128}

	Red25    = color.RGBA{255, 0, 0, 64}
	Green25  = color.RGBA{0, 255, 0, 64}
	Blue25   = color.RGBA{0, 0, 255, 64}
	Yellow25 = color.RGBA{255, 255, 0, 64}
	Cyan25   = color.RGBA{0, 255, 255, 64}
	Purple25 = color.RGBA{255, 0, 255, 64}
	White25  = color.RGBA{255, 255, 255, 64}

	// translate ids to name string
	ID = map[int]string{
		1: "player",
		2: "square",
		3: "wall",
		4: "tester",
		5: "finish",
		6: "checkpoint",
	}
)

// Constants
const (
	ScreenWidth  = 1150
	ScreenHeight = 864

	IDPlayer     = 1
	IDSquare     = 2
	IDWall       = 3
	IDTester     = 4
	IDFinish     = 5
	IDCheckpoint = 6
)
