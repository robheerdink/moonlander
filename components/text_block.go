package comp

import (
	"fmt"
	"log"
	"strconv"
	"time"

	con "moonlander/constants"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var (
	face     font.Face
	laps     Text
	laptime  Text
	gravity  Text
	endTimes string
)

// duration formater, stores duration as total MS, and seperate min, sec, ms
type duration struct {
	total, min, sec, ms int
}

// TextBlock is a positioal holder for multiple textblocks
type TextBlock struct {
	x, y int
}

func init() {
	// setup fonts
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	face = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// NewTextBlock contructor
func NewTextBlock(x, y int) TextBlock {
	// setup textboxes

	gravity = NewText(0, 0, "", face, con.Blue)
	laps = NewText(0, 30, "", face, con.Blue)
	laptime = NewText(0, 60, "", face, con.Red)
	return TextBlock{x, y}
}

// Draw it
func (o *TextBlock) Draw(screen *ebiten.Image) error {

	// update laps
	laps.text = "LAP: " + strconv.Itoa(LP.CurrentLap)

	// update lap laptime
	if !LP.LapStartTime.IsZero() {
		if LP.CurrentLap > LP.MaxLaps {
			if endTimes == "" {
				endTimes = calcEndTimes()
			}
			laptime.text = endTimes
		} else {
			et := getElapsedTime(LP.LapStartTime)
			laptime.text = fmt.Sprintf("%02d:%02d.%03d", et.min, et.sec, et.ms)
		}
	}

	// update Gravity value
	gravity.text = fmt.Sprintf("%.2fG", (WP.Gravity * 60))

	// draw
	text.Draw(screen, gravity.text, face, o.x+gravity.x, o.y+gravity.y, laps.color)
	text.Draw(screen, laps.text, face, o.x+laps.x, o.y+laps.y, laps.color)
	text.Draw(screen, laptime.text, face, o.x+laptime.x, o.y+laptime.y, laptime.color)
	return nil
}

// Get elapsed time from a start time
func getElapsedTime(t time.Time) duration {
	et := time.Now().Sub(t)
	return fmtDuration(et)
}

// fmt duration in total ms, min, sec, ms
func fmtDuration(et time.Duration) duration {
	total := int(et.Milliseconds())
	min := int(total/60000) % 60
	sec := int(total / 1000 % 60)
	ms := int(total % 1000)
	return duration{
		total: total, min: min, sec: sec, ms: ms,
	}
}

func calcEndTimes() string {
	var str string
	var tt time.Duration
	for _, lt := range LP.LapTimes {
		tt += lt
		lap := fmtDuration(lt)
		str += fmt.Sprintf("%02d:%02d.%03d\n", lap.min, lap.sec, lap.ms)
	}
	total := fmtDuration(tt)
	str += fmt.Sprintf("-----------\n%02d:%02d.%03d\n", total.min, total.sec, total.ms)
	return str
}

// atm we update values in draw
// func (o *TextBlock) Update(screen *ebiten.Image) error {
// }
