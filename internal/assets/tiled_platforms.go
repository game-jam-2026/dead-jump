package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

func CreateTiledPlatform(w *ecs.World, x, y float64, tilesWide int, tile *ebiten.Image) ecs.EntityID {
	entity := w.CreateEntity()

	tileW := float64(TileSize)
	tileH := float64(TileSize)
	platformW := float64(tilesWide) * tileW
	platformH := tileH

	platformImg := ebiten.NewImage(int(platformW), int(platformH))
	for i := 0; i < tilesWide; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i)*tileW, 0)
		platformImg.DrawImage(tile, op)
	}

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: platformImg,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, platformW, platformH),
	})
	w.SetComponent(entity, components.StaticBody())

	return entity
}

func CreateTiledPlatformTall(w *ecs.World, x, y float64, tilesWide, tilesHigh int, tile *ebiten.Image) ecs.EntityID {
	entity := w.CreateEntity()

	tileW := float64(TileSize)
	tileH := float64(TileSize)
	platformW := float64(tilesWide) * tileW
	platformH := float64(tilesHigh) * tileH

	// Создаем изображение платформы из тайлов
	platformImg := ebiten.NewImage(int(platformW), int(platformH))
	for j := 0; j < tilesHigh; j++ {
		for i := 0; i < tilesWide; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i)*tileW, float64(j)*tileH)
			platformImg.DrawImage(tile, op)
		}
	}

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image: platformImg,
	})
	w.SetComponent(entity, components.Collision{
		Shape: resolv.NewRectangleFromTopLeft(x, y, platformW, platformH),
	})
	w.SetComponent(entity, components.StaticBody())

	return entity
}

func CreateDecoration(w *ecs.World, x, y float64, tile *ebiten.Image) ecs.EntityID {
	entity := w.CreateEntity()

	w.SetComponent(entity, components.Position{
		Vector: linalg.Vector2{X: x, Y: y},
	})
	w.SetComponent(entity, components.Sprite{
		Image:  tile,
		ZIndex: -1,
	})

	return entity
}
