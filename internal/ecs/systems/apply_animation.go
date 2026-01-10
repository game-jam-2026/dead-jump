package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func ApplyAnimation(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Animation)(nil)).Elem(),
	)

	for _, entity := range entities {
		animation, err := ecs.GetComponent[components.Animation](world, entity)
		if err != nil {
			continue
		}

		img := animation.CheckAndGetImage()
		if img == nil {
			continue
		}

		world.SetComponent(entity, *animation)
		world.SetComponent(entity, components.Sprite{
			Image: img,
		})
	}
}
