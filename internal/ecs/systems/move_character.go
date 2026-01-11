package systems

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

const (
	MoveSpeed         = 0.5
	JumpForce         = 6.0
	StepSoundCooldown = 10
)

var stepSoundTimer int

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

	isMovingLeft := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	isMovingRight := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)

	if isMovingLeft {
		body.AddForce(linalg.Vector2{X: -MoveSpeed * body.Mass, Y: 0})
	}
	if isMovingRight {
		body.AddForce(linalg.Vector2{X: MoveSpeed * body.Mass, Y: 0})
	}

	if stepSoundTimer > 0 {
		stepSoundTimer--
	}
	if body.IsGrounded && (isMovingLeft || isMovingRight) && stepSoundTimer == 0 {
		audio.Play(audio.SoundStep)
		stepSoundTimer = StepSoundCooldown
	}

	if body.IsGrounded && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp)) {
		vel.Vector.Y = -JumpForce
		body.IsGrounded = false
		w.SetComponent(characterID, *vel)
	}

	w.SetComponent(characterID, *body)
}
