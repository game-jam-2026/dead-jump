package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/systems"
	"github.com/game-jam-2026/dead-jump/internal/levels"
)

type Game struct {
	w *ecs.World
}

func (g *Game) Update() error {
	systems.ApplyVelocity(g.w)
	systems.PushColliders(g.w)

	if key := g.isKeyJustPressed(); key != nil {
		systems.MoveCharacter(g.w, *key)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	systems.DrawSprites(g.w, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) isKeyJustPressed() *ebiten.Key {
	var key ebiten.Key
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		key = ebiten.KeySpace
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		key = ebiten.KeyRight
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		key = ebiten.KeyLeft
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		key = ebiten.KeyUp
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		key = ebiten.KeyDown
	default:
		return nil
	}
	return &key
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Game{
		w: levels.LoadLevel1(),
	}); err != nil {
		log.Fatal(err)
	}
}
