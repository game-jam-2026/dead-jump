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

func LoadLevel2() *ecs.World {
	w := ecs.NewWorld()
	CreateAudioManager(w)
	CreateLifeCounter(w, 3)
	CreateStartPoint(w, 20, 50)

	playerID := CreateCharacter(w, 10, 50, 1)

	CreateGround(w, 0, 210, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})
	CreateCannon(w, 32, 180, -math.Pi/2, 0)

	CreateExteriorObject(w, 220, 30, MoonImage)
	CreateExteriorObject(w, -10, 162, FirLeftImage)
	CreateExteriorObject(w, 128, 100, FirLeftImage)
	CreateExteriorObject(w, 160, 100, FirRightImage)

	CreateSpike(w, 64, 215, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})
	createWallBlock(w, 128, 180, 32, 320, 5, false, components.Repeatable{
		Direction: linalg.Vector2{Y: 1},
		Count:     2,
	})
	createWallBlock(w, 160, 180, 32, 320, 5, true, components.Repeatable{
		Direction: linalg.Vector2{Y: 1},
		Count:     2,
	})
	CreateGround(w, 128, 148, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})
	CreateSpike(w, 192, 215, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})
	CreateGround(w, 256, 210, 24, 24, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})
	CreateSpike(w, 288, 215, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     1,
	})

	createFloor(w, -20, 0, 20, 320)
	createFloor(w, 320, 0, 20, 320)
	CreateLevelFinish(w, 264, 190)

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

func createWallBlock(w *ecs.World, x, y, width, height float64, zIndex int, isRight bool,
	repeatable components.Repeatable) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := WallLeftImage
	if isRight {
		img = WallRightImage
	}

	w.SetComponent(entity, components.Sprite{
		Image:  img,
		ZIndex: zIndex,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(entity, components.StaticBody())
	w.SetComponent(entity, repeatable)
	ApplyRepeatable(w, entity)

	return entity
}
