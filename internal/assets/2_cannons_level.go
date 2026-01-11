package assets

import (
	"math"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func Load2CannonsLevel() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 5)

	CreateTiledPlatform(w, 0, 368, 31, TileGroundTextured)

	CreateTiledPlatform(w, 28, 336, 5, TileGroundGrass)

	if TileTree != nil {
		CreateDecoration(w, 16, 304, TileTree)
	}

	CreateStartPoint(w, 30, 320)
	playerID := CreateCharacter(w, 30, 320, 1.0)

	CreateSpike(w, 144, 352, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     8,
	})

	CreateTiledPlatform(w, 208, 336, 4, TileGroundGrass)

	cannonEntity := CreateCannon(w, 224, 320, -math.Pi*3/4)
	cannon, _ := ecs.GetComponent[components.Cannon](w, cannonEntity)
	cannon.BurstCount = 5
	cannon.BurstDelay = 7
	cannon.FireRate = 140
	cannon.ProjectileSpeed = 28.5
	cannon.ProjectileMass = 14.0
	w.SetComponent(cannonEntity, *cannon)

	CreateTiledPlatform(w, 32, 256, 5, TileGroundGrass)
	if TileTree != nil {
		CreateDecoration(w, 40, 224, TileTree)
	}

	cannonEntityMidPl := CreateCannon(w, 62, 240, -math.Pi/5)
	cannonMidPl, _ := ecs.GetComponent[components.Cannon](w, cannonEntityMidPl)
	cannonMidPl.BurstCount = 5
	cannonMidPl.BurstDelay = 7
	cannonMidPl.FireRate = 140
	cannonMidPl.FramesSinceLastShot = 70
	cannonMidPl.ProjectileSpeed = 28.5
	cannonMidPl.ProjectileMass = 14.0
	w.SetComponent(cannonEntityMidPl, *cannonMidPl)

	CreateTiledPlatform(w, 224, 192, 7, TileGroundGrass)

	CreateSpike(w, 304, 176, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreateTiledPlatformTall(w, 368, 96, 8, 2, TileGroundTextured)
	CreateTiledPlatform(w, 368, 80, 8, TileGroundGrass)

	if TileTree != nil {
		CreateDecoration(w, 376, 48, TileTree)
		CreateDecoration(w, 416, 48, TileTree)
	}

	CreateLevelFinish(w, 432, 55)

	if TileColumn != nil {
		for i := 0; i < 23; i++ {
			CreateDecoration(w, 0, float64(i*16), TileColumn)
		}
	}

	if TileColumn != nil {
		for i := 0; i < 23; i++ {
			CreateDecoration(w, 484, float64(i*16), TileColumn)
		}
	}

	CreateGround(w, 0, 0, 8, 400, components.Repeatable{})
	CreateGround(w, 492, 0, 8, 400, components.Repeatable{})
	CreateGround(w, 0, 0, 500, 8, components.Repeatable{})

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.SetBounds(0, 0, 500, 384)
	camera.Smoothing = 0.08
	camera.DeadZoneX = 70
	camera.DeadZoneY = 30
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}
