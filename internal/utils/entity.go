package utils

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func KillEntity(
	w *ecs.World,
	entity ecs.EntityID,
	deadImage *ebiten.Image,
	scale float64,
	createCharacterFunc func(w *ecs.World, x, y float64, scale float64) ecs.EntityID,
) {
	audio.Play(audio.SoundDeath)

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

	w.SetComponent(entity, components.StaticBody())

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

	newPosY := pos.Vector.Y + 10

	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(pos.Vector.X, newPosY, width, height),
	})
	newVec := linalg.Vector2{
		X: pos.Vector.X,
		Y: newPosY,
	}
	w.SetComponent(entity, components.Position{Vector: newVec})

	startPoints := w.GetEntities(reflect.TypeOf((*components.StartPoint)(nil)).Elem())
	if len(startPoints) > 0 {
		spPos, _ := ecs.GetComponent[components.Position](w, startPoints[0])
		createCharacterFunc(w, spPos.Vector.X, spPos.Vector.Y, scale)
	}

	lifeCounters := w.GetEntities(reflect.TypeOf((*components.Life)(nil)).Elem())
	if len(lifeCounters) > 0 {
		life, _ := ecs.GetComponent[components.Life](w, lifeCounters[0])
		life.Count += -1
		w.SetComponent(lifeCounters[0], *life)
	}
}
