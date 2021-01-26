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
	Name = map[int]string{
		0: "unknown",
		1: "player",
		2: "background",
		3: "square",
		4: "wall",
		5: "tester",
		6: "finish",
		7: "checkpoint",
	}
)

// Constants
const (
	ScreenWidth  = 1280 //40 * 32
	ScreenHeight = 960  //30 * 32

	IDPlayer     = 1
	IDBG         = 2
	IDSquare     = 3
	IDWall       = 4
	IDTester     = 5
	IDFinish     = 6
	IDCheckpoint = 7
)
