package assets

import (
	"image/color"
	"math"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

func LoadCannonLevel() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)

	CreateStartPoint(w, 50, 280)
	playerID := CreateCharacter(w, 50, 280, 0.5)

	createFloor(w, 0, 320, 180, 16)

	CreateCannon(w, 20, 304, -math.Pi/4)

	CreateSpike(w, 100, 296, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     3,
	})

	createWallBlock(w, 200, 170, 20, 170)

	createFloor(w, 220, 320, 180, 16)

	createFloor(w, 396, 320, 100, 16)

	createWallBlock(w, 0, 0, 8, 400)
	createWallBlock(w, 492, 0, 8, 400)
	createWallBlock(w, 0, 392, 500, 8)

	camera := components.NewCamera(320, 240)
	camera.Target = int64(playerID)
	camera.SetBounds(0, 0, 500, 400)
	camera.Smoothing = 0.1
	camera.DeadZoneX = 20
	camera.DeadZoneY = 15
	w.SetResource(camera)

	w.SetResource(*physics.DefaultConfig())

	return w
}

func createFloor(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{80, 80, 80, 255})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(entity, components.StaticBody())

	return entity
}

func createWallBlock(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{60, 60, 60, 255})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(entity, components.StaticBody())

	return entity
}
