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

	heroImg, _, _ := image.Decode(bytes.NewReader(heroPNG))

	characterEntity1 := w.CreateEntity()
	w.SetComponent(characterEntity1, components.Position{
		Vector: linalg.Vector2{X: 100, Y: 100},
	})
	w.SetComponent(characterEntity1, components.Sprite{
		Image: GenerateRandomImage(24, 24),
	})
	w.SetComponent(characterEntity1, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(100, 100, 24, 24),
	})

	characterEntity2 := w.CreateEntity()
	w.SetComponent(characterEntity2, components.Position{
		Vector: linalg.Vector2{X: 110, Y: 50},
	})
	w.SetComponent(characterEntity2, components.Sprite{
		Image: ebiten.NewImageFromImage(heroImg),
	})
	w.SetComponent(characterEntity2, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(110, 50, 24, 32),
	})
	w.SetComponent(characterEntity2, components.Velocity{
		Vector: linalg.Vector2{X: 0, Y: 0.2},
	})

	return w
}
