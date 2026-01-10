package components

import "github.com/game-jam-2026/dead-jump/pkg/linalg"

type Camera struct {
	Position       linalg.Vector2
	ViewportWidth  float64
	ViewportHeight float64
	Target         int64
	Smoothing      float64
	DeadZoneX      float64
	DeadZoneY      float64
	MinX, MaxX     float64
	MinY, MaxY     float64
}

func NewCamera(viewportWidth, viewportHeight float64) Camera {
	return Camera{
		Position:       linalg.Zero(),
		ViewportWidth:  viewportWidth,
		ViewportHeight: viewportHeight,
		Target:         0,
		Smoothing:      0.15,
		DeadZoneX:      10,
		DeadZoneY:      10,
	}
}

func (c *Camera) SetBounds(minX, minY, maxX, maxY float64) {
	c.MinX = minX
	c.MinY = minY
	c.MaxX = maxX
	c.MaxY = maxY
}

func (c *Camera) WorldToScreen(worldPos linalg.Vector2) linalg.Vector2 {
	return linalg.Vector2{
		X: worldPos.X - c.Position.X,
		Y: worldPos.Y - c.Position.Y,
	}
}

func (c *Camera) IsVisible(worldPos linalg.Vector2, width, height float64) bool {
	return worldPos.X+width > c.Position.X &&
		worldPos.X < c.Position.X+c.ViewportWidth &&
		worldPos.Y+height > c.Position.Y &&
		worldPos.Y < c.Position.Y+c.ViewportHeight
}
