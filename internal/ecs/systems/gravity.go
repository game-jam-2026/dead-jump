package systems

import (
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
)

func ApplyGravity(world *ecs.World, cfg *physics.Config) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil {
			continue
		}

		if body.IsKinematic || body.GravityScale == 0 {
			continue
		}

		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		gravityForce := cfg.Gravity.Scale(body.GravityScale)
		vel.Vector = vel.Vector.Add(gravityForce)

		if vel.Vector.Y > cfg.TerminalVelocity {
			vel.Vector.Y = cfg.TerminalVelocity
		}

		if body.MaxSpeed > 0 {
			vel.Vector = vel.Vector.ClampLength(body.MaxSpeed)
		}

		world.SetComponent(e, *vel)
	}
}

func ApplySlopeGravity(world *ecs.World, cfg *physics.Config) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil {
			continue
		}

		if !body.IsGrounded || body.IsKinematic {
			continue
		}

		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		surface := findContactSurface(world, e)
		if surface == nil || surface.SlopeAngle == 0 {
			continue
		}

		frictionFactor := 1.0 - surface.Friction
		if frictionFactor < 0.1 {
			frictionFactor = 0.1
		}

		slideForce := surface.GetSlideAcceleration(cfg.Gravity.Y) * body.GravityScale * frictionFactor
		vel.Vector.X -= slideForce

		world.SetComponent(e, *vel)
	}
}
