package components

import "github.com/hajimehoshi/ebiten/v2"

type Character struct {
	GroundedSprite *ebiten.Image
	JumpingSprite  *ebiten.Image
}
