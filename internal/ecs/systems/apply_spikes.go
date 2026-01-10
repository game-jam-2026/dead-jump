package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func ApplySpikes(world *ecs.World) {
	characters := world.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
	)

	spikes := world.GetEntities(
		reflect.TypeOf((*components.Spike)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
	)

	for _, charEntity := range characters {
		charCollision, err := ecs.GetComponent[components.Collision](world, charEntity)
		if err != nil {
			continue
		}

		for _, spikeEntity := range spikes {
			spikeCollision, err := ecs.GetComponent[components.Collision](world, spikeEntity)
			if err != nil {
				continue
			}

			intersection := charCollision.Shape.Intersection(spikeCollision.Shape)
			if !intersection.IsEmpty() {
				assets.KillEntity(world, charEntity)
				break
			}
		}
	}
}
