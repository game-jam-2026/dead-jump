package assets

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

var (
	//go:embed img/hero.png
	HeroPNG []byte
	//go:embed img/hero_jump.png
	HeroJumpPNG []byte
	//go:embed img/dead_hero.png
	DeadHeroPNG []byte
	//go:embed img/spike.png
	SpikePNG []byte
	//go:embed img/ground.png
	GroundPNG []byte
	//go:embed img/orange.png
	OrangePNG []byte
	//go:embed img/wall_left.png
	WallLeftPNG []byte
	//go:embed img/wall_right.png
	WallRightPNG []byte
	//go:embed img/fir_right.png
	FirRightPNG []byte
	//go:embed img/fir_left.png
	FirLeftPNG []byte
	//go:embed img/moon.png
	MoonPNG []byte
	//go:embed img/tombstone1.png
	Tombstone1PNG []byte
	//go:embed img/tombstone2.png
	Tombstone2PNG []byte
	//go:embed img/tombstone3.png
	Tombstone3PNG []byte
	//go:embed img/cannon_left.png
	CannonLeftPNG []byte
	//go:embed img/cannon_top.png
	CannonTopPNG []byte
	//go:embed img/cannon_right.png
	CannonRightPNG []byte
)

var (
	HeroImage        *ebiten.Image
	HeroJumpImage    *ebiten.Image
	DeadHeroImage    *ebiten.Image
	SpikeImage       *ebiten.Image
	GroundImage      *ebiten.Image
	OrangeImage      *ebiten.Image
	Tombstone1Image  *ebiten.Image
	Tombstone2Image  *ebiten.Image
	Tombstone3Image  *ebiten.Image
	WallLeftImage    *ebiten.Image
	WallRightImage   *ebiten.Image
	FirRightImage    *ebiten.Image
	FirLeftImage     *ebiten.Image
	MoonImage        *ebiten.Image
	CannonLeftImage  *ebiten.Image
	CannonTopImage   *ebiten.Image
	CannonRightImage *ebiten.Image
)

func init() {
	heroImg, _, _ := image.Decode(bytes.NewReader(HeroPNG))
	HeroImage = ebiten.NewImageFromImage(heroImg)

	heroJumpImg, _, _ := image.Decode(bytes.NewReader(HeroJumpPNG))
	HeroJumpImage = ebiten.NewImageFromImage(heroJumpImg)

	deadHeroImg, _, _ := image.Decode(bytes.NewReader(DeadHeroPNG))
	DeadHeroImage = ebiten.NewImageFromImage(deadHeroImg)

	spikeImg, _, _ := image.Decode(bytes.NewReader(SpikePNG))
	SpikeImage = ebiten.NewImageFromImage(spikeImg)

	groundImg, _, _ := image.Decode(bytes.NewReader(GroundPNG))
	GroundImage = ebiten.NewImageFromImage(groundImg)

	orangeImg, _, _ := image.Decode(bytes.NewReader(OrangePNG))
	OrangeImage = ebiten.NewImageFromImage(orangeImg)

	wallLeftImg, _, _ := image.Decode(bytes.NewReader(WallLeftPNG))
	WallLeftImage = ebiten.NewImageFromImage(wallLeftImg)

	wallRightImg, _, _ := image.Decode(bytes.NewReader(WallRightPNG))
	WallRightImage = ebiten.NewImageFromImage(wallRightImg)

	firRightImg, _, _ := image.Decode(bytes.NewReader(FirRightPNG))
	FirRightImage = ebiten.NewImageFromImage(firRightImg)

	firLeftImg, _, _ := image.Decode(bytes.NewReader(FirLeftPNG))
	FirLeftImage = ebiten.NewImageFromImage(firLeftImg)

	moonImg, _, _ := image.Decode(bytes.NewReader(MoonPNG))
	MoonImage = ebiten.NewImageFromImage(moonImg)

	tombstone1Img, _, _ := image.Decode(bytes.NewReader(Tombstone1PNG))
	Tombstone1Image = ebiten.NewImageFromImage(tombstone1Img)

	tombstone2Img, _, _ := image.Decode(bytes.NewReader(Tombstone2PNG))
	Tombstone2Image = ebiten.NewImageFromImage(tombstone2Img)

	tombstone3Img, _, _ := image.Decode(bytes.NewReader(Tombstone3PNG))
	Tombstone3Image = ebiten.NewImageFromImage(tombstone3Img)

	cannonLeftImg, _, _ := image.Decode(bytes.NewReader(CannonLeftPNG))
	CannonLeftImage = ebiten.NewImageFromImage(cannonLeftImg)

	cannonTopImg, _, _ := image.Decode(bytes.NewReader(CannonTopPNG))
	CannonTopImage = ebiten.NewImageFromImage(cannonTopImg)

	cannonRightImg, _, _ := image.Decode(bytes.NewReader(CannonRightPNG))
	CannonRightImage = ebiten.NewImageFromImage(cannonRightImg)
}

func CreateCharacter(w *ecs.World, x, y float64, scale float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := HeroImage.Bounds()
	origW := float64(bounds.Dx())
	origH := float64(bounds.Dy())
	width := origW * scale
	height := origH * scale

	groundedSprite := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	groundedSprite.DrawImage(HeroImage, op)

	// jumping sprite
	jumpBounds := HeroJumpImage.Bounds()
	jumpOrigW := float64(jumpBounds.Dx())
	jumpOrigH := float64(jumpBounds.Dy())
	jumpWidth := jumpOrigW * scale
	jumpHeight := jumpOrigH * scale

	jumpingSprite := ebiten.NewImage(int(jumpWidth), int(jumpHeight))
	opJump := &ebiten.DrawImageOptions{}
	opJump.GeoM.Scale(scale, scale)
	jumpingSprite.DrawImage(HeroJumpImage, opJump)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: groundedSprite,
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

	w.SetComponent(entity, components.Character{
		GroundedSprite: groundedSprite,
		JumpingSprite:  jumpingSprite,
	})

	return entity
}

func CreateLifeCounter(w *ecs.World, lifeCnt int) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: 10, Y: 10},
	})
	w.SetComponent(entity, components.Sprite{
		Image: OrangeImage,
	})
	w.SetComponent(entity, components.Life{
		Count: lifeCnt,
	})
	w.SetComponent(entity, components.Repeatable{
		Direction: linalg.Vector2{X: 1},
		Count:     lifeCnt,
	})
	w.SetComponent(entity, components.ScreenSpace{})

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

func CreateExteriorObject(w *ecs.World, x, y float64, img *ebiten.Image) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: img,
	})
	w.SetComponent(entity, components.StaticBody())

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

func CreatePlatform(w *ecs.World, x, y, width, height float64, repeatable components.Repeatable) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := GroundImage.Bounds()
	halfHeight := bounds.Dy() / 2
	halfSprite := ebiten.NewImage(bounds.Dx(), halfHeight)
	op := &ebiten.DrawImageOptions{}
	halfSprite.DrawImage(GroundImage.SubImage(image.Rect(0, 0, bounds.Dx(), halfHeight)).(*ebiten.Image), op)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: halfSprite,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y+8, width, height/2),
	})
	w.SetComponent(entity, components.StaticBody())
	w.SetComponent(entity, repeatable)

	ApplyRepeatable(w, entity)

	return entity
}

func CreateTombstone1(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: Tombstone1Image,
	})
	return entity
}

func CreateTombstone2(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: Tombstone2Image,
	})
	return entity
}

func CreateTombstone3(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()
	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: Tombstone3Image,
	})
	return entity
}

func CreateAudioManager(w *ecs.World) ecs.EntityID {
	entity := w.CreateEntity()

	audio.Init()

	w.SetComponent(entity, components.AudioContext{Context: audio.GetContext()})
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

func CreateCannon(w *ecs.World, x, y float64, direction float64, position int) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})

	var img *ebiten.Image
	switch {
	case position < 0:
		img = CannonLeftImage
	case position > 0:
		img = CannonRightImage
	default:
		img = CannonTopImage
	}

	w.SetComponent(entity, components.Sprite{
		Image:  img,
		ZIndex: 3,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, 32, 32),
	})
	w.SetComponent(entity, components.StaticBody())

	cannon := components.DefaultCannon()
	cannon.Direction = direction
	cannon.FireRate = 120
	cannon.ProjectileSpeed = 15.0
	cannon.ProjectileMass = 15.0
	w.SetComponent(entity, cannon)

	return entity
}

func CreateLevelFinish(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := OrangeImage.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	tintedImg := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(0.3, 0.5, 1.5, 1.0)
	tintedImg.DrawImage(OrangeImage, op)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image:  tintedImg,
		ZIndex: 5,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})
	w.SetComponent(entity, components.LevelFinish{})

	return entity
}

func CreateCorpse(w *ecs.World, x, y float64, scale float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := DeadHeroImage.Bounds()
	width := float64(bounds.Dx()) * scale
	height := float64(bounds.Dy()) * scale

	scaledImg := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(DeadHeroImage, op)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: scaledImg,
	})
	w.SetComponent(entity, components.Corpse{
		Durability: -1,
		IsSettled:  true,
	})
	w.SetComponent(entity, components.StaticBody())

	return entity
}

func CreateEpilogueFinish(w *ecs.World, x, y float64) ecs.EntityID {
	entity := w.CreateEntity()

	bounds := OrangeImage.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	tintedImg := ebiten.NewImage(int(width), int(height))
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(1.5, 1.0, 0.3, 1.0)
	tintedImg.DrawImage(OrangeImage, op)

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image:  tintedImg,
		ZIndex: 5,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, width, height),
	})
	w.SetComponent(entity, components.EpilogueFinish{})

	return entity
}
