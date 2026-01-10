package systems

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawSprites(world *ecs.World, screen *ebiten.Image) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Sprite)(nil)).Elem(),
	)

	// чтоб не мигало при пересечении, но вообще – надо бы какой-то порядок отрисовки придумать
	sort.Slice(entities, func(i, j int) bool {
		return entities[i] > entities[j]
	})

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
		op.GeoM.Translate(pos.Vector.X, pos.Vector.Y)
		screen.DrawImage(sprite.Image, op)
	}
}
