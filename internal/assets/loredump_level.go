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

const loreText = "You are a mage in a dying world. Your life is supported by the magic fruits, even if your body dies, making you to experience the same places over and over again. Try to find out what has happened by collecting warp fruits to traverse the world."

func LoadLoreDumpLevel() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 3)

	levelWidth := 320.0
	levelHeight := 240.0

	CreateStartPoint(w, 30, 150)
	playerID := CreateCharacter(w, 30, 150, 0.5)

	CreateGround(w, 0, 210, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     14,
	})

	createLoreWall(w, 0, 0, 8, levelHeight)
	createLoreWall(w, levelWidth-8, 0, 8, levelHeight)

	CreateLevelFinish(w, levelWidth-50, 190)

	w.SetResource(components.LoreText{Text: loreText})

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

func createLoreWall(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{30, 25, 40, 255})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(entity, components.StaticBody())

	return entity
}
