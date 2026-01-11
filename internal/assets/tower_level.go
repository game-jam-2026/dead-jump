package assets

import (
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

const TowerLevelHeight = 600

func LoadTowerLevel() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 3)
	CreateStartPoint(w, 40, 530)

	playerID := CreateCharacter(w, 40, 530, 0.5)

	CreateSpike(w, 0, 584, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     20,
	})

	CreatePlatform(w, 40, 550, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 120, 530, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 200, 510, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 260, 490, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 200, 470, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 160, 450, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 100, 430, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 20, 410, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 100, 390, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	CreatePlatform(w, 180, 370, 24, 16, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	// Tombstones on some platforms
	CreateTombstone1(w, 124, 500)
	CreateTombstone2(w, 204, 440)
	CreateTombstone3(w, 24, 380)

	CreateLevelFinish(w, 274, 354)

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.MinX = 0
	camera.MaxX = 320
	camera.MinY = 0
	camera.MaxY = TowerLevelHeight
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}
