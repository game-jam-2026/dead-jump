package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func DrawLifeCounter(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Life)(nil)).Elem(),
	)
	if len(entities) == 0 {
		return
	}

	entity := entities[0]

	life, err := ecs.GetComponent[components.Life](world, entity)
	if err != nil {
		return
	}

	world.SetComponent(entity, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     life.Count,
	})
	world.SetComponent(entity, components.Sprite{
		Image: assets.HeartImage,
	})

	if life.Count == 0 {
		world.DestroyEntity(entity)
	}

	assets.ApplyRepeatable(world, entity)
}
