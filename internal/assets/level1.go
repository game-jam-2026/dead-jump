package assets

import (
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func LoadLevel1() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 5)
	CreateStartPoint(w, 20, 50)

	playerID := CreateCharacter(w, 20, 50, 1)

	CreateGround(w, 0, 210, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})

	CreateSpike(w, 64, 220, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     6,
	})

	CreateGround(w, 256, 210, 32, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     5,
	})

	CreateLevelFinish(w, 290, 194)

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.SetBounds(0, 0, WorldWidth, WorldHeight)
	camera.Smoothing = 0.1
	camera.DeadZoneX = 20
	camera.DeadZoneY = 15
	camera.MinX = 0
	camera.MaxY = 240
	camera.MaxX = 320
	camera.MinY = 0
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}
