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

	CreateWall(w, 50, 100, 24, 24)
	CreateStartPoint(w, 110, 50)
	playerID := CreateCharacter(w, 110, 50, 0.5)

	CreateWall(w, 50, 100, 24, 24)
	CreateSpike(w, 100, 180, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     7,
	})

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.SetBounds(0, 0, WorldWidth, WorldHeight)
	camera.Smoothing = 0.1
	camera.DeadZoneX = 20
	camera.DeadZoneY = 15
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}
