package systems

import (
	"image/color"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawCollisions(world *ecs.World, screen *ebiten.Image) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
	)

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}
		collision, err := ecs.GetComponent[components.Collision](world, e)
		if err != nil {
			continue
		}

		bounds := collision.Shape.Bounds()
		width := float32(bounds.Width())
		height := float32(bounds.Height())

		clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		vector.StrokeRect(screen, float32(pos.Vector.X), float32(pos.Vector.Y), width, height, 1, clr, false)
	}
}
