package levels

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

//go:embed static/hero.png
var heroPNG []byte

func GenerateRandomImage(width, height int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 255,
			}
			img.Set(x, y, c)
		}
	}

	return ebiten.NewImageFromImage(img)
}

func LoadLevel1() *ecs.World {
	w := ecs.NewWorld()
	makeFloor(w)
	makeCharacter(w)

	return w
}

func makeFloor(w *ecs.World) {
	floor := w.CreateEntity()
	w.SetComponent(floor, components.Position{
		Vector: linalg.Vector2{X: 0, Y: 200},
	})
	w.SetComponent(floor, components.Sprite{
		Image: GenerateRandomImage(320, 40),
	})
	w.SetComponent(floor, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(0, 200, 320, 40),
	})
}

func makeCharacter(w *ecs.World) {
	img, _, _ := image.Decode(bytes.NewReader(heroPNG))
	character := w.CreateEntity()

	w.SetComponent(character, components.Position{
		Vector: linalg.Vector2{X: 110, Y: 10},
	})
	w.SetComponent(character, components.Sprite{
		Image: ebiten.NewImageFromImage(img),
	})
	w.SetComponent(character, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(110, 10, 24, 32),
	})
	w.SetComponent(character, components.Velocity{
		Vector: linalg.Vector2{X: 0, Y: 1},
	})
	w.SetComponent(character, components.Character{})
}
