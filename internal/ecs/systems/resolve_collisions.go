package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func PushColliders(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	for i, e := range entities {
		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		c, err := ecs.GetComponent[components.Collision](world, e)
		if err != nil {
			continue
		}

		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}

		bounds := c.Shape.Bounds()
		c.Shape.SetPosition(pos.Vector.X+bounds.Width()/2, pos.Vector.Y+bounds.Height()/2)

		for _, otherE := range entities[i+1:] {
			otherC, err := ecs.GetComponent[components.Collision](world, otherE)
			if err != nil {
				continue
			}

			intersection := c.Shape.Intersection(otherC.Shape)
			if intersection.IsEmpty() {
				continue
			}

			pos.Vector.X += intersection.MTV.X
			pos.Vector.Y += intersection.MTV.Y

			if intersection.MTV.X != 0 {
				vel.Vector.X = 0
			}
			if intersection.MTV.Y != 0 {
				vel.Vector.Y = 0
			}

			c.Shape.SetPosition(pos.Vector.X+bounds.Width()/2, pos.Vector.Y+bounds.Height()/2)
			world.SetComponent(e, *pos)
			world.SetComponent(e, *c)
			world.SetComponent(e, *vel)
			break
		}
	}
}
