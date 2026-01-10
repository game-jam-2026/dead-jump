package main

import (
	"log"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/systems"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	w *ecs.World
}

func (g *Game) Update() error {
	systems.ApplyVelocity(g.w)
	systems.PushColliders(g.w)
	systems.ApplySpikes(g.w)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	systems.DrawSprites(g.w, screen)
	systems.DrawCollisions(g.w, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Game{
		w: assets.LoadLevel1(),
	}); err != nil {
		log.Fatal(err)
	}
}
