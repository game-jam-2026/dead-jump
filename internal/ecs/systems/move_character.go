package systems

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

const (
	MoveSpeed = 0.5
	JumpForce = 6.0
)

func MoveCharacter(w *ecs.World) {
	entities := w.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
		reflect.TypeOf((*components.PhysicsBody)(nil)).Elem(),
		reflect.TypeOf((*components.Velocity)(nil)).Elem(),
	)

	if len(entities) == 0 {
		return
	}

	characterID := entities[0]

	body, err := ecs.GetComponent[components.PhysicsBody](w, characterID)
	if err != nil {
		return
	}

	vel, err := ecs.GetComponent[components.Velocity](w, characterID)
	if err != nil {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		body.AddForce(linalg.Vector2{X: -MoveSpeed * body.Mass, Y: 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		body.AddForce(linalg.Vector2{X: MoveSpeed * body.Mass, Y: 0})
	}

	if body.IsGrounded && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp)) {
		vel.Vector.Y = -JumpForce
		body.IsGrounded = false
		w.SetComponent(characterID, *vel)
	}

	w.SetComponent(characterID, *body)
}
