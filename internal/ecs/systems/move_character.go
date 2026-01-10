package systems

import (
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func MoveCharacter(w *ecs.World, key ebiten.Key) {
	const accelerationCoefficient = 0.5
	entities := w.GetEntities(reflect.TypeOf((*components.Character)(nil)).Elem())
	if len(entities) == 0 {
		return
	}

	character := entities[0]
	velocity, err := ecs.GetComponent[components.Velocity](w, character)
	if err != nil {
		fmt.Println(err)
	}

	switch key {
	case ebiten.KeySpace, ebiten.KeyUp:
		if velocity.Vector.Y > -accelerationCoefficient {
			velocity.Vector = velocity.Vector.Add(linalg.Vector2{X: 0, Y: -accelerationCoefficient})
		}
	case ebiten.KeyRight:
		if velocity.Vector.X < accelerationCoefficient {
			velocity.Vector = velocity.Vector.Add(linalg.Vector2{X: accelerationCoefficient, Y: 0})
		}
	case ebiten.KeyLeft:
		if velocity.Vector.X > -accelerationCoefficient {
			velocity.Vector = velocity.Vector.Add(linalg.Vector2{X: -accelerationCoefficient, Y: 0})
		}
	case ebiten.KeyDown:
		if velocity.Vector.Y < accelerationCoefficient {
			velocity.Vector = velocity.Vector.Add(linalg.Vector2{X: 0, Y: accelerationCoefficient})
		}
	}

	w.SetComponent(character, *velocity)
}
