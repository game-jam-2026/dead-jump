package utils

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func KillEntity(
	w *ecs.World,
	entity ecs.EntityID,
	deadImage *ebiten.Image,
	scale float64,
	createCharacterFunc func(w *ecs.World, x, y float64, scale float64) ecs.EntityID,
) {
	assets.PlayRandomDeathSound()

	pos, _ := ecs.GetComponent[components.Position](w, entity)

	err := w.RemoveComponent(entity, components.Character{})
	if err != nil {
		panic(err)
	}
	err = w.RemoveComponent(entity, components.Velocity{})
	if err != nil {
		panic(err)
	}
	_ = w.RemoveComponent(entity, components.PhysicsBody{})

	w.SetComponent(entity, components.Corpse{
		Durability: 5,
	})

	bounds := deadImage.Bounds()
	width := float64(bounds.Dx()) * scale
	height := float64(bounds.Dy()) * scale

	scaledImg := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(deadImage, op)

	w.SetComponent(entity, components.Sprite{
		Image: scaledImg,
	})

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(pos.Vector.X, pos.Vector.Y, width, height),
	})

	startPoints := w.GetEntities(reflect.TypeOf((*components.StartPoint)(nil)).Elem())
	if len(startPoints) > 0 {
		spPos, _ := ecs.GetComponent[components.Position](w, startPoints[0])
		createCharacterFunc(w, spPos.Vector.X, spPos.Vector.Y, scale)
	}
}
