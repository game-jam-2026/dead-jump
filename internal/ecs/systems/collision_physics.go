package systems

import (
	"math"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

type CollisionResult struct {
	EntityA     ecs.EntityID
	EntityB     ecs.EntityID
	MTV         linalg.Vector2
	Normal      linalg.Vector2
	Penetration float64
}

func ResolvePhysicsCollisions(world *ecs.World, cfg *physics.Config) []CollisionResult {
	var results []CollisionResult

	entities := world.GetEntities(
		reflect.TypeOf((*components.Collision)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	resetGroundedState(world)
	syncCollisionPositions(world, entities)

	for i := 0; i < cfg.CollisionIterations; i++ {
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
	}

	return results
}

func resolveCollision(
	world *ecs.World, cfg *physics.Config,
	entityA, entityB ecs.EntityID,
	posA, posB *components.Position,
	colA, colB *components.Collision,
	bodyA, bodyB *components.PhysicsBody,
	velA, velB *components.Velocity,
	mtv, normal linalg.Vector2,
) {
	massA := float64(1)
	massB := float64(1)
	staticA := false
	staticB := false

	if bodyA != nil {
		if bodyA.IsStatic() {
			staticA = true
			massA = math.MaxFloat64
		} else {
			massA = bodyA.Mass
		}
	}

	if bodyB != nil {
		if bodyB.IsStatic() {
			staticB = true
			massB = math.MaxFloat64
		} else {
			massB = bodyB.Mass
		}
	}

	if staticA && staticB {
		return
	}

	totalMass := massA + massB
	ratioA := massB / totalMass
	ratioB := massA / totalMass

	if staticA {
		ratioA = 0
		ratioB = 1
	}
	if staticB {
		ratioA = 1
		ratioB = 0
	}

	posA.Vector = posA.Vector.Add(mtv.Scale(ratioA))
	posB.Vector = posB.Vector.Sub(mtv.Scale(ratioB))

	boundsA := colA.Shape.Bounds()
	boundsB := colB.Shape.Bounds()
	colA.Shape.SetPosition(posA.Vector.X+boundsA.Width()/2, posA.Vector.Y+boundsA.Height()/2)
	colB.Shape.SetPosition(posB.Vector.X+boundsB.Width()/2, posB.Vector.Y+boundsB.Height()/2)

	world.SetComponent(entityA, *posA)
	world.SetComponent(entityB, *posB)
	world.SetComponent(entityA, *colA)
	world.SetComponent(entityB, *colB)

	if velA != nil || velB != nil {
		resolveVelocities(world, entityA, entityB, bodyA, bodyB, velA, velB, normal)
	}

	updateGroundedState(world, entityA, entityB, bodyA, bodyB, normal)
}

func resolveVelocities(
	world *ecs.World,
	entityA, entityB ecs.EntityID,
	bodyA, bodyB *components.PhysicsBody,
	velA, velB *components.Velocity,
	normal linalg.Vector2,
) {
	va := linalg.Zero()
	vb := linalg.Zero()
	if velA != nil {
		va = velA.Vector
	}
	if velB != nil {
		vb = velB.Vector
	}

	relativeVel := va.Sub(vb)
	velAlongNormal := relativeVel.Dot(normal)

	if velAlongNormal > 0 {
		return
	}

	restitution := float64(0)
	if bodyA != nil {
		restitution = math.Max(restitution, bodyA.Bounciness)
	}
	if bodyB != nil {
		restitution = math.Max(restitution, bodyB.Bounciness)
	}

	surfaceA := findContactSurface(world, entityA)
	surfaceB := findContactSurface(world, entityB)
	if surfaceA != nil && surfaceA.Bounciness > restitution {
		restitution = surfaceA.Bounciness
	}
	if surfaceB != nil && surfaceB.Bounciness > restitution {
		restitution = surfaceB.Bounciness
	}

	minBounceVelocity := 1.0
	if math.Abs(velAlongNormal) < minBounceVelocity {
		restitution = 0
	}

	invMassA := float64(0)
	invMassB := float64(0)

	if bodyA != nil && !bodyA.IsStatic() && bodyA.Mass > 0 {
		invMassA = 1 / bodyA.Mass
	}
	if bodyB != nil && !bodyB.IsStatic() && bodyB.Mass > 0 {
		invMassB = 1 / bodyB.Mass
	}

	if invMassA == 0 && invMassB == 0 {
		return
	}

	j := -(1 + restitution) * velAlongNormal
	j /= invMassA + invMassB

	impulse := normal.Scale(j)

	if velA != nil && invMassA > 0 {
		velA.Vector = velA.Vector.Add(impulse.Scale(invMassA))
		if normal.Y < -0.5 && velA.Vector.Y > 0 {
			velA.Vector.Y = 0
		}
		world.SetComponent(entityA, *velA)
	}

	if velB != nil && invMassB > 0 {
		velB.Vector = velB.Vector.Sub(impulse.Scale(invMassB))
		if normal.Y > 0.5 && velB.Vector.Y > 0 {
			velB.Vector.Y = 0
		}
		world.SetComponent(entityB, *velB)
	}
}

func ApplyProjectileImpulse(world *ecs.World, projectile, target ecs.EntityID, impulseMagnitude float64) {
	projPos, err := ecs.GetComponent[components.Position](world, projectile)
	if err != nil {
		return
	}

	targetBody, err := ecs.GetComponent[components.PhysicsBody](world, target)
	if err != nil || targetBody.IsStatic() {
		return
	}

	targetVel, err := ecs.GetComponent[components.Velocity](world, target)
	if err != nil {
		targetVel = &components.Velocity{Vector: linalg.Zero()}
	}

	targetPos, err := ecs.GetComponent[components.Position](world, target)
	if err != nil {
		return
	}

	impactDir := targetPos.Vector.Sub(projPos.Vector).Normalized()
	if impactDir.IsZero() {
		impactDir = linalg.Up()
	}

	impulseVel := targetBody.AddImpulse(impactDir.Scale(impulseMagnitude))
	targetVel.Vector = targetVel.Vector.Add(impulseVel)

	if targetBody.IsGrounded {
		targetVel.Vector.Y -= impulseMagnitude * 0.3
		targetBody.IsGrounded = false
	}

	world.SetComponent(target, *targetVel)
	world.SetComponent(target, *targetBody)
}

func resetGroundedState(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
	)

	for _, e := range entities {
		body, err := ecs.GetComponent[components.PhysicsBody](world, e)
		if err != nil || body.IsKinematic {
			continue
		}

		body.IsGrounded = false
		body.GroundNormal = linalg.Up()
		world.SetComponent(e, *body)
	}
}

func updateGroundedState(world *ecs.World, entityA, entityB ecs.EntityID, bodyA, bodyB *components.PhysicsBody, normal linalg.Vector2) {
	if normal.Y < -0.5 && bodyA != nil && !bodyA.IsKinematic {
		bodyA.IsGrounded = true
		bodyA.GroundNormal = normal.Scale(-1)
		world.SetComponent(entityA, *bodyA)
	}

	if normal.Y > 0.5 && bodyB != nil && !bodyB.IsKinematic {
		bodyB.IsGrounded = true
		bodyB.GroundNormal = normal
		world.SetComponent(entityB, *bodyB)
	}
}

func syncCollisionPositions(world *ecs.World, entities []ecs.EntityID) {
	for _, e := range entities {
		col, err := ecs.GetComponent[components.Collision](world, e)
		if err != nil {
			continue
		}

		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}

		bounds := col.Shape.Bounds()
		col.Shape.SetPosition(pos.Vector.X+bounds.Width()/2, pos.Vector.Y+bounds.Height()/2)
		world.SetComponent(e, *col)
	}
}
