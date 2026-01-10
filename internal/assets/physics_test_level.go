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

const (
	WorldWidth  = 800
	WorldHeight = 400
)

func LoadPhysicsTestLevel() *ecs.World {
	w := ecs.NewWorld()

	CreateAudioManager(w)

	CreateStartPoint(w, 50, 300)
	playerID := CreateCharacter(w, 50, 300, 0.5)

	createPlatform(w, 0, 350, 100, 16, components.SurfaceNormal, color.RGBA{100, 100, 100, 255})
	createPlatform(w, 100, 350, 80, 16, components.SurfaceIce, color.RGBA{200, 220, 255, 255})
	createPlatform(w, 180, 350, 80, 16, components.SurfaceRough, color.RGBA{139, 90, 43, 255})
	createPlatform(w, 260, 350, 60, 16, components.SurfaceBouncy, color.RGBA{255, 100, 150, 255})

	CreateSpike(w, 130, 334, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     2,
	})

	createPlatform(w, 320, 350, 230, 16, components.SurfaceNormal, color.RGBA{90, 90, 90, 255})

	CreateSpike(w, 480, 334, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     3,
	})

	createSlopeVisual(w, 340, 280, 80, 70, math.Pi/12, color.RGBA{150, 120, 80, 255})
	createSlopeVisual(w, 430, 250, 70, 100, math.Pi/6, color.RGBA{140, 100, 70, 255})

	createPlatform(w, 350, 170, 40, 8, components.SurfaceNormal, color.RGBA{80, 80, 80, 255})

	createPlatform(w, 550, 350, 250, 16, components.SurfaceNormal, color.RGBA{100, 100, 100, 255})

	createSlopeVisual(w, 570, 320, 40, 30, math.Pi/8, color.RGBA{120, 100, 90, 255})
	createSlopeVisual(w, 620, 280, 40, 40, math.Pi/8, color.RGBA{120, 100, 90, 255})

	createPlatform(w, 700, 150, 80, 8, components.SurfaceNormal, color.RGBA{80, 80, 80, 255})

	createWall(w, 0, 0, 8, WorldHeight)
	createWall(w, WorldWidth-8, 0, 8, WorldHeight)
	createWall(w, 0, WorldHeight-8, WorldWidth, 8)

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

func createSlopeVisual(w *ecs.World, x, y, width, height float64, angle float64, col color.RGBA) ecs.EntityID {
	return createSlopeWithSurface(w, x, y, width, height, angle, components.SurfaceNormal, col)
}

func createSlopeWithSurface(w *ecs.World, x, y, width, height float64, angle float64, surfaceType components.SurfaceType, col color.RGBA) ecs.EntityID {
	slope := w.CreateEntity()

	w.SetComponent(slope, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	for px := 0; px < int(width); px++ {
		for py := 0; py < int(height); py++ {
			fx := float64(px)
			fy := float64(py)

			threshold := height * (1.0 - fx/width)

			if fy >= threshold {
				depth := (fy - threshold) / height
				darken := uint8(depth * 30)

				edgeDist := fy - threshold
				alpha := uint8(255)
				if edgeDist < 1.0 {
					alpha = uint8(edgeDist * 255)
				}

				c := color.RGBA{
					R: col.R - darken,
					G: col.G - darken,
					B: col.B - darken,
					A: alpha,
				}
				img.Set(px, py, c)
			}
		}
	}
	w.SetComponent(slope, components.Sprite{
		Image: img,
	})

	polygon := resolv.NewConvexPolygon(
		x+width/2, y+height/2,
		[]float64{
			-width / 2, height / 2,
			width / 2, height / 2,
			width / 2, -height / 2,
		},
	)

	w.SetComponent(slope, components.Collision{
		Shape: polygon,
	})

	w.SetComponent(slope, components.StaticBody())

	surface := components.NewSlopedSurface(surfaceType, angle)
	w.SetComponent(slope, surface)

	return slope
}

func createPlatform(w *ecs.World, x, y, width, height float64, surfaceType components.SurfaceType, col color.RGBA) ecs.EntityID {
	platform := w.CreateEntity()

	w.SetComponent(platform, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(col)
	w.SetComponent(platform, components.Sprite{
		Image: img,
	})

	w.SetComponent(platform, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(platform, components.StaticBody())

	w.SetComponent(platform, components.NewSurface(surfaceType))

	return platform
}

func createWall(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	wall := w.CreateEntity()

	w.SetComponent(wall, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{50, 50, 50, 255})
	w.SetComponent(wall, components.Sprite{
		Image: img,
	})

	w.SetComponent(wall, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	w.SetComponent(wall, components.StaticBody())

	return wall
}
