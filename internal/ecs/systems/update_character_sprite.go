package systems

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func UpdateCharacterSprite(w *ecs.World) {
	entities := w.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Sprite)(nil)).Elem(),
	)

	if len(entities) == 0 {
		return
	}

	characterID := entities[0]

	character, err := ecs.GetComponent[components.Character](w, characterID)
	if err != nil {
		return
	}

	body, err := ecs.GetComponent[components.PhysicsBody](w, characterID)
	if err != nil {
		return
	}

	sprite, err := ecs.GetComponent[components.Sprite](w, characterID)
	if err != nil {
		return
	}

	if character.GroundedSprite == nil || character.JumpingSprite == nil {
		return
	}

	var targetSprite *ebiten.Image
	if body.IsGrounded {
		targetSprite = character.GroundedSprite
	} else {
		targetSprite = character.JumpingSprite
	}

	if sprite.Image != targetSprite {
		sprite.Image = targetSprite
		w.SetComponent(characterID, *sprite)
	}
}
