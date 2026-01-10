package systems

import (
	"math"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
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

		pos.Vector = pos.Vector.Add(vel.Vector)

		world.SetComponent(e, *pos)
	}
}

func ApplyVelocityWithCollisions(world *ecs.World, cfg *physics.Config) []CollisionResult {
	var allResults []CollisionResult

	maxSpeed := findMaxSpeed(world)

	substeps := 1
	if cfg.MaxStepDistance > 0 && maxSpeed > cfg.MaxStepDistance {
		substeps = int(math.Ceil(maxSpeed / cfg.MaxStepDistance))
	}

	if substeps > 16 {
		substeps = 16
	}

	resetGroundedState(world)

	for step := 0; step < substeps; step++ {
		applyPartialVelocity(world, substeps)
		results := resolveCollisionsSubstep(world, cfg)
		allResults = append(allResults, results...)
	}

	return allResults
}

func findMaxSpeed(world *ecs.World) float64 {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	maxSpeed := 0.0
	for _, e := range entities {
		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		speed := vel.Vector.Length()
		if speed > maxSpeed {
			maxSpeed = speed
		}
	}

	return maxSpeed
}

func applyPartialVelocity(world *ecs.World, substeps int) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	fraction := 1.0 / float64(substeps)

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}
		vel, err := ecs.GetComponent[components.Velocity](world, e)
		if err != nil {
			continue
		}

		pos.Vector = pos.Vector.Add(vel.Vector.Scale(fraction))

		world.SetComponent(e, *pos)
	}
}

func resolveCollisionsSubstep(world *ecs.World, cfg *physics.Config) []CollisionResult {
	var results []CollisionResult

	entities := world.GetEntities(
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	syncCollisionPositions(world, entities)

	for j, entityA := range entities {
		colA, err := ecs.GetComponent[components.Collision](world, entityA)
		if err != nil {
			continue
		}

		posA, err := ecs.GetComponent[components.Position](world, entityA)
		if err != nil {
			continue
		}

		for _, entityB := range entities[j+1:] {
			colB, err := ecs.GetComponent[components.Collision](world, entityB)
			if err != nil {
				continue
			}

			intersection := colA.Shape.Intersection(colB.Shape)
			if intersection.IsEmpty() {
				continue
			}

			posB, err := ecs.GetComponent[components.Position](world, entityB)
			if err != nil {
				continue
			}

			bodyA, _ := ecs.GetComponent[components.PhysicsBody](world, entityA)
			bodyB, _ := ecs.GetComponent[components.PhysicsBody](world, entityB)

			velA, _ := ecs.GetComponent[components.Velocity](world, entityA)
			velB, _ := ecs.GetComponent[components.Velocity](world, entityB)

			mtv := linalg.Vector2{X: intersection.MTV.X, Y: intersection.MTV.Y}
			normal := mtv.Normalized()

			result := CollisionResult{
				EntityA:     entityA,
				EntityB:     entityB,
				MTV:         mtv,
				Normal:      normal,
				Penetration: mtv.Length(),
			}
			results = append(results, result)

			resolveCollision(
				world, cfg,
				entityA, entityB,
				posA, posB,
				colA, colB,
				bodyA, bodyB,
				velA, velB,
				mtv, normal,
			)
		}
	}

	syncCollisionPositions(world, entities)

	return results
}
