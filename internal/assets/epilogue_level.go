package assets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

const epilogueLoreText = "You remember this place... Need one more fruit to break the cycle."

func LoadEpilogueLevel() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 1)

	levelWidth := 320.0
	levelHeight := 240.0

	CreateStartPoint(w, 30, 180)
	playerID := CreateCharacter(w, 30, 180, 1)

	CreateGround(w, 0, 210, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     14,
	})

	createEpilogueWall(w, 0, 0, 8, levelHeight)
	createEpilogueWall(w, levelWidth-8, 0, 8, levelHeight)

	CreateCorpse(w, 50, 202, 1)
	CreateCorpse(w, 70, 202, 1)
	CreateCorpse(w, 55, 197, 1)

	CreateCorpse(w, 110, 202, 1)
	CreateCorpse(w, 125, 200, 1)

	CreateCorpse(w, 150, 202, 1)
	CreateCorpse(w, 165, 198, 1)
	CreateCorpse(w, 175, 202, 1)

	CreateCorpse(w, 210, 202, 1)
	CreateCorpse(w, 225, 200, 1)
	CreateCorpse(w, 240, 202, 1)
	CreateCorpse(w, 255, 198, 1)

	CreateTombstone1(w, 90, 182)
	CreateTombstone2(w, 185, 182)
	CreateTombstone3(w, 270, 182)

	CreateEpilogueFinish(w, levelWidth-50, 190)

	w.SetResource(components.LoreText{Text: epilogueLoreText})

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.SetBounds(0, 0, levelWidth, levelHeight)
	camera.Smoothing = 0.1
	camera.DeadZoneX = 20
	camera.DeadZoneY = 15
	camera.MinX = 0
	camera.MinY = 0
	camera.MaxX = 320
	camera.MaxY = 240
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}

func createEpilogueWall(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{20, 15, 25, 255})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(entity, components.StaticBody())

	return entity
}
