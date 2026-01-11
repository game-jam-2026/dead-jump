package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func ApplyLevelFinish(world *ecs.World) bool {
	characters := world.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
	)

	finishTriggers := world.GetEntities(
		reflect.TypeOf((*components.LevelFinish)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
	)

	for _, charEntity := range characters {
		charCollision, err := ecs.GetComponent[components.Collision](world, charEntity)
		if err != nil {
			continue
		}

		for _, finishEntity := range finishTriggers {
			finishCollision, err := ecs.GetComponent[components.Collision](world, finishEntity)
			if err != nil {
				continue
			}

			intersection := charCollision.Shape.Intersection(finishCollision.Shape)
			if !intersection.IsEmpty() {
				return true
			}
		}
	}
	return false
}
