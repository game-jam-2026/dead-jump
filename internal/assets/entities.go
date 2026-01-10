package assets

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

//go:embed static/hero.png
var HeroPNG []byte

//go:embed static/dead_hero.png
var DeadHeroPNG []byte

//go:embed static/spike.png
var SpikePNG []byte

var (
	HeroImage     *ebiten.Image
	DeadHeroImage *ebiten.Image
	SpikeImage    *ebiten.Image
)

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

func init() {
	heroImg, _, _ := image.Decode(bytes.NewReader(HeroPNG))
	HeroImage = ebiten.NewImageFromImage(heroImg)

	deadHeroImg, _, _ := image.Decode(bytes.NewReader(DeadHeroPNG))
	DeadHeroImage = ebiten.NewImageFromImage(deadHeroImg)

	spikeImg, _, _ := image.Decode(bytes.NewReader(SpikePNG))
	SpikeImage = ebiten.NewImageFromImage(spikeImg)
}

func CreateCharacter(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: HeroImage,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, 24, 32),
	})
	w.SetComponent(entity, components.Velocity{
		Vector: linalg.Vector2{X: 0, Y: 0.5},
	})
	w.SetComponent(entity, components.Character{})

	return entity
}

func CreateSpike(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := SpikeImage.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: SpikeImage,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})
	w.SetComponent(entity, components.Spike{})

	return entity
}

func CreateStartPoint(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.StartPoint{})

	return entity
}

func CreateWall(w *ecs.World, x, y, width, height float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: GenerateRandomImage(int(width), int(height)),
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})

	return entity
}

func KillEntity(w *ecs.World, entity ecs.EntityID) {
	pos, _ := ecs.GetComponent[components.Position](w, entity)

	err := w.RemoveComponent(entity, components.Character{})
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
		Image: DeadHeroImage,
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
		CreateCharacter(w, spPos.Vector.X, spPos.Vector.Y)
	}
}
