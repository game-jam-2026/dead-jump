package assets

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
)

//go:embed img/hero.png
var HeroPNG []byte

//go:embed img/dead_hero.png
var DeadHeroPNG []byte

//go:embed img/spike.png
var SpikePNG []byte

//go:embed img/ground.png
var GroundPNG []byte

var (
	HeroImage     *ebiten.Image
	DeadHeroImage *ebiten.Image
	SpikeImage    *ebiten.Image
	GroundImage   *ebiten.Image
)

func init() {
	heroImg, _, _ := image.Decode(bytes.NewReader(HeroPNG))
	HeroImage = ebiten.NewImageFromImage(heroImg)

	deadHeroImg, _, _ := image.Decode(bytes.NewReader(DeadHeroPNG))
	DeadHeroImage = ebiten.NewImageFromImage(deadHeroImg)

	spikeImg, _, _ := image.Decode(bytes.NewReader(SpikePNG))
	SpikeImage = ebiten.NewImageFromImage(spikeImg)

	groundImg, _, _ := image.Decode(bytes.NewReader(GroundPNG))
	GroundImage = ebiten.NewImageFromImage(groundImg)
}

func CreateCharacter(w *ecs.World, x, y float64, scale float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := HeroImage.Bounds()
	origW := float64(bounds.Dx())
	origH := float64(bounds.Dy())
	width := origW * scale
	height := origH * scale

	scaledImg := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(HeroImage, op)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: scaledImg,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})
	w.SetComponent(entity, components.Velocity{
		Vector: linalg.Zero(),
	})

	body := components.DefaultPhysicsBody()
	body.Mass = 1.0
	body.Friction = 0.3
	body.AirDrag = 0.15
	body.GravityScale = 1.0
	body.MaxSpeed = 20.0
	w.SetComponent(entity, body)

	w.SetComponent(entity, components.Character{})

	return entity
}

func CreateSpike(w *ecs.World, x, y float64, repeat components.Repeatable) ecs.EntityID {
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
	w.SetComponent(entity, repeat)
	w.SetComponent(entity, components.StaticBody())

	ApplyRepeatable(w, entity)

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

func CreateGround(w *ecs.World, x, y, width, height float64, repeatable components.Repeatable) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: GroundImage,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y+16, width, height),
	})
	w.SetComponent(entity, components.StaticBody())
	w.SetComponent(entity, repeatable)

	ApplyRepeatable(w, entity)

	return entity
}

func CreateAudioManager(w *ecs.World) ecs.EntityID {
	audioContext := audio.NewContext(44100)
	InitAudio(audioContext)
	entity := w.CreateEntity()
	w.SetComponent(entity, components.AudioContext{Context: audioContext})
	return entity
}

func ApplyRepeatable(w *ecs.World, entity ecs.EntityID) {
	rep, err := ecs.GetComponent[components.Repeatable](w, entity)
	if err != nil {
		return
	}

	sprite, err := ecs.GetComponent[components.Sprite](w, entity)
	if err != nil {
		fmt.Println(err)
		return
	}
	pos, err := ecs.GetComponent[components.Position](w, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	origBounds := sprite.Image.Bounds()
	origW, origH := origBounds.Dx(), origBounds.Dy()

	newW := origW
	newH := origH
	if rep.Direction.X > 0 {
		newW = origW * rep.Count
	}
	if rep.Direction.Y > 0 {
		newH = origH * rep.Count
	}

	newImg := ebiten.NewImage(newW, newH)
	for i := 0; i < rep.Count; i++ {
		op := &ebiten.DrawImageOptions{}
		offsetX := float64(i*origW) * math.Abs(rep.Direction.X)
		offsetY := float64(i*origH) * math.Abs(rep.Direction.Y)

		op.GeoM.Translate(offsetX, offsetY)
		newImg.DrawImage(sprite.Image, op)
	}

	w.SetComponent(entity, components.Sprite{Image: newImg})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(pos.Vector.X, pos.Vector.Y, float64(newW), float64(newH)),
	})
}

func CreateCannon(w *ecs.World, x, y float64, direction float64) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	size := 16
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{40, 40, 40, 255})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})

	cannon := components.DefaultCannon()
	cannon.Direction = direction
	cannon.FireRate = 120
	cannon.ProjectileSpeed = 15.0
	cannon.ProjectileMass = 15.0
	w.SetComponent(entity, cannon)

	return entity
}
