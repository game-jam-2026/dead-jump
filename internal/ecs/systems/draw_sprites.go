package systems

import (
	"fmt"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawSprites(world *ecs.World, screen *ebiten.Image) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Sprite)(nil)).Elem(),
	)

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sprite, err := ecs.GetComponent[components.Sprite](world, e)
		if err != nil {
			fmt.Println(err)
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.P.X, pos.P.Y)
		screen.DrawImage(sprite.Image, op)
	}
}
