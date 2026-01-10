package utils

import (
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func KillEntity(
	w *ecs.World,
	entity ecs.EntityID,
	deadImage *ebiten.Image,
	deathSounds [][]byte,
	createCharacterFunc func(w *ecs.World, x, y float64, scale float64) ecs.EntityID,
) {
	err := PlayRandomSound(w, deathSounds)
	if err != nil {
		fmt.Println(err)
	}

	pos, _ := ecs.GetComponent[components.Position](w, entity)

	err = w.RemoveComponent(entity, components.Character{})
	if err != nil {
		panic(err)
	}
	err = w.RemoveComponent(entity, components.Velocity{})
	if err != nil {
		panic(err)
	}

	w.SetComponent(entity, components.Corpse{
		Durability: 5,
	})

	w.SetComponent(entity, components.Sprite{
		Image: deadImage,
	})

	// Оффсет для насаживания на штык
	newY := pos.Vector.Y + 12
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(pos.Vector.X, newY, 24, 24),
	})
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: pos.Vector.X, Y: newY},
	})

	startPoints := w.GetEntities(reflect.TypeOf((*components.StartPoint)(nil)).Elem())
	if len(startPoints) > 0 {
		spPos, _ := ecs.GetComponent[components.Position](w, startPoints[0])
		createCharacterFunc(w, spPos.Vector.X, spPos.Vector.Y, 0.5)
	}
}
