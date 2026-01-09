package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func ApplyVelocity(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}
		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		pos.P = pos.P.Add(vel.V)

		world.SetComponent(e, *pos)
	}
}
