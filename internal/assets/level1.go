package assets

import (
	"github.com/game-jam-2026/dead-jump/internal/ecs"
)

func LoadLevel1() *ecs.World {
	w := ecs.NewWorld()

	CreateWall(w, 50, 100, 24, 24)
	CreateStartPoint(w, 110, 50)
	CreateCharacter(w, 110, 50)
	CreateSpike(w, 100, 180)

	return w
}
