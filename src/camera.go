package src

import (
	"fmt"
	"math"
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort    f64.Vec2
	Position    f64.Vec2
	ZoomFactor  int
	Rotation    int
	boundRight  float64
	boundBottom float64
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

// Render camera to screen
func (c *Camera) Render(world, screen *ebiten.Image) error {
	return screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

// ScreenToWorld calc screen position to world position
func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happend that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

// Reset Camera to default values
func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
	c.boundRight = float64(sha.LP.Width) - float64(sha.ScreenWidth)
	c.boundBottom = float64(sha.LP.Height) - float64(sha.ScreenHeight)
	fmt.Println(sha.LP.Width, sha.LP.Height, sha.ScreenWidth, sha.ScreenHeight)
}

// Update ..
func (c *Camera) Update() error {
	panCamera := false

	// pan WSAD
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.Position[0] -= 5
		panCamera = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.Position[0] += 5
		panCamera = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Position[1] -= 5
		panCamera = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Position[1] += 5
		panCamera = true
	}

	// center camera view on player, unless panning
	if !panCamera {
		x, y, _, _ := player.GetImageInfo()
		c.Position[0] = x - sha.ScreenWidth/2
		c.Position[1] = y - sha.ScreenHeight/2

		//fmt.Println(x, y)
		fmt.Println(c, panCamera)
	}

	// set level bounds
	if c.Position[0] < 0 {
		c.Position[0] = 0
	}
	if c.Position[0] > c.boundRight {
		c.Position[0] = c.boundRight
	}
	if c.Position[1] < 0 {
		c.Position[1] = 0
	}
	if c.Position[1] > c.boundBottom {
		c.Position[1] = c.boundBottom
	}

	// rotate Q/W
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		c.Rotation--
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.Rotation++
	}

	// zoom SHIFT/CTRL
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		if c.ZoomFactor < 100 {
			c.ZoomFactor++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if c.ZoomFactor > -100 {
			c.ZoomFactor--
		}
	}

	// reset SPACE
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Reset()
	}

	return nil
}
