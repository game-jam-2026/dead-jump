package assets

import (
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func LoadLevel1() *ecs.World {
	w := ecs.NewWorld()

	CreateWall(w, 50, 100, 24, 24)
	CreateStartPoint(w, 110, 50)
	CreateCharacter(w, 110, 50)
	CreateSpike(w, 100, 180, components.Repeatable{
		Direction: linalg.Vector2{
			X: 1,
		},
		Count: 7,
	})

	return w
}
