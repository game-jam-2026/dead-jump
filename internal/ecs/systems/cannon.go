package systems

import (
	"image/color"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func UpdateCannons(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Cannon)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	for _, e := range entities {
		cannon, err := ecs.GetComponent[components.Cannon](world, e)
		if err != nil || !cannon.Active {
			continue
		}

		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}

		cannon.FramesSinceLastShot++

		if cannon.FramesSinceLastShot >= cannon.FireRate {
			cannon.FramesSinceLastShot = 0

			spawnDistance := float64(12)
			cannonCenterX := pos.Vector.X + 8
			cannonCenterY := pos.Vector.Y + 4

			spawnX := cannonCenterX + math.Cos(cannon.Direction)*spawnDistance
			spawnY := cannonCenterY + math.Sin(cannon.Direction)*spawnDistance

			velocity := linalg.Vector2{
				X: math.Cos(cannon.Direction) * cannon.ProjectileSpeed,
				Y: math.Sin(cannon.Direction) * cannon.ProjectileSpeed,
			}

			spawnProjectile(world, spawnX, spawnY, velocity, cannon.ProjectileMass)
			audio.Play(audio.SoundCannonShot)
		}

		world.SetComponent(e, *cannon)
	}
}

func spawnProjectile(world *ecs.World, x, y float64, velocity linalg.Vector2, mass float64) ecs.EntityID {
	projectile := world.CreateEntity()

	world.SetComponent(projectile, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	size := 8
	img := ebiten.NewImage(size, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			dx := float64(i) - float64(size)/2
			dy := float64(j) - float64(size)/2
			if dx*dx+dy*dy < float64(size*size)/4 {
				img.Set(i, j, color.RGBA{50, 50, 50, 255})
			}
		}
	}
	world.SetComponent(projectile, components.Sprite{Image: img})

	world.SetComponent(projectile, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, float64(size), float64(size)),
	})

	world.SetComponent(projectile, components.Velocity{Vector: velocity})

	body := components.ProjectileBody(mass)
	body.Bounciness = 0.0
	world.SetComponent(projectile, body)

	world.SetComponent(projectile, components.Projectile{
		ImpulseMagnitude:   mass * 3,
		DestroyOnHit:       false,
		Lifetime:           300,
		MinSpeedForImpulse: 0.5,
	})

	return projectile
}

func HandleProjectileCollisions(world *ecs.World, collisions []CollisionResult) {
	for _, col := range collisions {
		projA, errA := ecs.GetComponent[components.Projectile](world, col.EntityA)
		projB, errB := ecs.GetComponent[components.Projectile](world, col.EntityB)

		var projectileID ecs.EntityID
		var targetID ecs.EntityID
		var proj *components.Projectile

		if errA == nil {
			projectileID = col.EntityA
			targetID = col.EntityB
			proj = projA
		} else if errB == nil {
			projectileID = col.EntityB
			targetID = col.EntityA
			proj = projB
		} else {
			continue
		}

		if proj.IsStationary {
			continue
		}

		body, err := ecs.GetComponent[components.PhysicsBody](world, targetID)
		if err == nil && body.IsStatic() {
			world.DestroyEntity(projectileID)
			continue
		}

		_, isCharacter := ecs.GetComponent[components.Character](world, targetID)
		if isCharacter == nil {
			audio.Play(audio.SoundProjectileHit)
			projVel, err := ecs.GetComponent[components.Velocity](world, projectileID)
			if err == nil && projVel.Vector.Length() >= proj.MinSpeedForImpulse {
				ApplyProjectileImpulse(world, projectileID, targetID, proj.ImpulseMagnitude)
				projVel.Vector = projVel.Vector.Scale(0.1)
				world.SetComponent(projectileID, *projVel)
			}

			if proj.DestroyOnHit {
				world.DestroyEntity(projectileID)
			}
		}
	}
}

func UpdateProjectileLifetime(world *ecs.World) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Projectile)(nil)).Elem(),
	)

	stationaryThreshold := 0.5

	for _, e := range entities {
		proj, err := ecs.GetComponent[components.Projectile](world, e)
		if err != nil {
			continue
		}

		if !proj.IsStationary {
			vel, velErr := ecs.GetComponent[components.Velocity](world, e)
			if velErr == nil && vel.Vector.Length() < stationaryThreshold {
				proj.IsStationary = true

				body, bodyErr := ecs.GetComponent[components.PhysicsBody](world, e)
				if bodyErr == nil {
					body.GravityScale = 0
					body.IsKinematic = true
					world.SetComponent(e, *body)
				}

				vel.Vector.X = 0
				vel.Vector.Y = 0
				world.SetComponent(e, *vel)
			}
		}

		if proj.Lifetime < 0 {
			world.SetComponent(e, *proj)
			continue
		}

		proj.Lifetime--

		if proj.Lifetime <= 0 {
			world.DestroyEntity(e)
		} else {
			world.SetComponent(e, *proj)
		}
	}
}

func CleanupOffscreenProjectiles(world *ecs.World, screenWidth, screenHeight float64) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Projectile)(nil)).Elem(),
		reflect.TypeOf((*components.Position)(nil)).Elem(),
	)

	margin := float64(50)

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}

		if pos.Vector.X < -margin || pos.Vector.X > screenWidth+margin ||
			pos.Vector.Y < -margin || pos.Vector.Y > screenHeight+margin {
			world.DestroyEntity(e)
		}
	}
}
