package systems

import (
	"math"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func ApplyFriction(world *ecs.World, cfg *physics.Config) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil {
			continue
		}

		if body.IsKinematic {
			continue
		}

		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		if body.IsGrounded {
			friction := body.Friction

			surfaceFriction := cfg.DefaultFriction
			contactSurface := findContactSurface(world, e)
			if contactSurface != nil {
				surfaceFriction = contactSurface.Friction
			}

			effectiveFriction := math.Sqrt(friction * surfaceFriction)
			frictionFactor := 1.0 - effectiveFriction

			vel.Vector.X *= frictionFactor

			if math.Abs(vel.Vector.X) < cfg.MinVelocity {
				vel.Vector.X = 0
			}
		} else {
			if body.AirDrag > 0 {
				dragFactor := 1.0 - body.AirDrag
				vel.Vector.X *= dragFactor
			}
		}

		world.SetComponent(e, *vel)
	}
}

func ApplyConveyorBelt(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil || !body.IsGrounded || body.IsKinematic {
			continue
		}

		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		surface := findContactSurface(world, e)
		if surface == nil || surface.Type != components.SurfaceConveyor {
			continue
		}

		conveyorVel := surface.ConveyorDirection.Scale(surface.ConveyorSpeed)
		vel.Vector = vel.Vector.Add(conveyorVel)

		world.SetComponent(e, *vel)
	}
}

func ApplyAccumulated(world *ecs.World, cfg *physics.Config) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil {
			continue
		}

		if body.IsKinematic {
			continue
		}

		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		vel.Vector = vel.Vector.Add(body.Acceleration)

		if body.MaxSpeed > 0 {
			vel.Vector = vel.Vector.ClampLength(body.MaxSpeed)
		}

		body.Acceleration = linalg.Zero()

		world.SetComponent(e, *vel)
		world.SetComponent(e, *body)
	}
}

func findContactSurface(world *ecs.World, entity ecs.EntityID) *components.Surface {
	col, err := ecs.GetComponent[components.Collision](world, entity)
	if err != nil {
		return nil
	}

	pos, err := ecs.GetComponent[components.Position](world, entity)
	if err != nil {
		return nil
	}

	entityBounds := col.Shape.Bounds()
	entityBottom := pos.Vector.Y + entityBounds.Height()

	surfaceEntities := world.GetEntities(
		reflect.TypeOf((*components.Surface)(nil)).Elem(),
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	for _, se := range surfaceEntities {
		if se == entity {
			continue
		}

		surfacePos, err := ecs.GetComponent[components.Position](world, se)
		if err != nil {
			continue
		}

		surfaceCol, err := ecs.GetComponent[components.Collision](world, se)
		if err != nil {
			continue
		}

		surfaceBounds := surfaceCol.Shape.Bounds()

		entityLeft := pos.Vector.X
		entityRight := pos.Vector.X + entityBounds.Width()
		surfaceLeft := surfacePos.Vector.X
		surfaceRight := surfacePos.Vector.X + surfaceBounds.Width()

		horizontalOverlap := entityRight > surfaceLeft && entityLeft < surfaceRight

		surfaceTop := surfacePos.Vector.Y
		verticalContact := entityBottom >= surfaceTop-2 && entityBottom <= surfaceTop+8

		if horizontalOverlap && verticalContact {
			surface, err := ecs.GetComponent[components.Surface](world, se)
			if err == nil {
				return surface
			}
		}
	}

	return nil
}
