package main

import (
	"log"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/ecs/systems"
	"github.com/game-jam-2026/dead-jump/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	w *ecs.World
}

func NewGame() *Game {
	return &Game{
		w: assets.LoadPhysicsTestLevel(),
	}
}

func (g *Game) Update() error {
	systems.MoveCharacter(g.w)
	systems.UpdateCannons(g.w)

	cfg, _ := ecs.GetResource[physics.Config](g.w)

	systems.ApplyGravity(g.w, cfg)
	systems.ApplyAccumulated(g.w, cfg)
	collisions := systems.ApplyVelocityWithCollisions(g.w, cfg)
	systems.HandleProjectileCollisions(g.w, collisions)
	systems.ApplySpikes(g.w)
	systems.ApplyAnimation(g.w)
	systems.ApplySlopeGravity(g.w, cfg)
	systems.ApplyFriction(g.w, cfg)
	systems.ApplyConveyorBelt(g.w)
	systems.UpdateProjectileLifetime(g.w)
	systems.CleanupOffscreenProjectiles(g.w, assets.WorldWidth, assets.WorldHeight)

	g.updateCameraTarget()
	systems.UpdateCameraSystem(g.w)

	return nil
}

func (g *Game) updateCameraTarget() {
	camera, err := ecs.GetResource[components.Camera](g.w)
	if err != nil {
		return
	}

	entities := g.w.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
	)
	if len(entities) > 0 {
		camera.Target = int64(entities[0])
		g.w.SetResource(*camera)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	camera, _ := ecs.GetResource[components.Camera](g.w)
	systems.DrawSpritesWithCamera(g.w, screen, camera)
	systems.DrawCollisions(g.w, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Dead Jump")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
