package systems

import (
	"math"
	"reflect"
	"sort"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"

	"github.com/hajimehoshi/ebiten/v2"
)

func UpdateCameraSystem(world *ecs.World) {
	camera, err := ecs.GetResource[components.Camera](world)
	if err != nil || camera.Target == 0 {
		return
	}

	updateCameraPosition(world, camera)
	world.SetResource(*camera)
}

func UpdateCamera(world *ecs.World, camera *components.Camera) {
	if camera.Target == 0 {
		return
	}

	updateCameraPosition(world, camera)
}

func updateCameraPosition(world *ecs.World, camera *components.Camera) {

	targetPos, err := ecs.GetComponent[components.Position](world, ecs.EntityID(camera.Target))
	if err != nil {
		return
	}

	var targetWidth, targetHeight float64 = 16, 16
	if col, err := ecs.GetComponent[components.Collision](world, ecs.EntityID(camera.Target)); err == nil {
		bounds := col.Shape.Bounds()
		targetWidth = bounds.Width()
		targetHeight = bounds.Height()
	}

	targetCenterX := targetPos.Vector.X + targetWidth/2
	targetCenterY := targetPos.Vector.Y + targetHeight/2

	desiredX := targetCenterX - camera.ViewportWidth/2
	desiredY := targetCenterY - camera.ViewportHeight/2

	if camera.MaxX > camera.MinX {
		desiredX = clamp(desiredX, camera.MinX, camera.MaxX-camera.ViewportWidth)
	}
	if camera.MaxY > camera.MinY {
		desiredY = clamp(desiredY, camera.MinY, camera.MaxY-camera.ViewportHeight)
	}

	smoothFactorX := 0.08
	smoothFactorY := 0.05

	camera.Position.X = lerp(camera.Position.X, desiredX, smoothFactorX)
	camera.Position.Y = lerp(camera.Position.Y, desiredY, smoothFactorY)

	camera.Position.X = math.Round(camera.Position.X)
	camera.Position.Y = math.Round(camera.Position.Y)
}

func DrawSpritesWithCamera(world *ecs.World, screen *ebiten.Image, camera *components.Camera) {
	entities := world.GetEntities(
		reflect.TypeOf((*components.Position)(nil)).Elem(),
		reflect.TypeOf((*components.Sprite)(nil)).Elem(),
	)

	sort.Slice(entities, func(i, j int) bool {
		spriteI, _ := ecs.GetComponent[components.Sprite](world, entities[i])
		spriteJ, _ := ecs.GetComponent[components.Sprite](world, entities[j])

		if spriteI.ZIndex != spriteJ.ZIndex {
			return spriteI.ZIndex < spriteJ.ZIndex
		}

		posI, _ := ecs.GetComponent[components.Position](world, entities[i])
		posJ, _ := ecs.GetComponent[components.Position](world, entities[j])

		if posI.Vector.Y != posJ.Vector.Y {
			return posI.Vector.Y < posJ.Vector.Y
		}

		return entities[i] < entities[j]
	})

	for _, e := range entities {
		pos, err := ecs.GetComponent[components.Position](world, e)
		if err != nil {
			continue
		}

		sprite, err := ecs.GetComponent[components.Sprite](world, e)
		if err != nil {
			continue
		}

		bounds := sprite.Image.Bounds()
		spriteWidth := float64(bounds.Dx())
		spriteHeight := float64(bounds.Dy())

		_, isScreenSpace := ecs.GetComponent[components.ScreenSpace](world, e)

		var screenPos linalg.Vector2
		if isScreenSpace == nil {
			screenPos = pos.Vector
		} else {
			if !camera.IsVisible(pos.Vector, spriteWidth, spriteHeight) {
				continue
			}
			screenPos = camera.WorldToScreen(pos.Vector)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenPos.X, screenPos.Y)
		screen.DrawImage(sprite.Image, op)
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
