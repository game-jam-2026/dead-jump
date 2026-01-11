package assets

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed img/tile_ground_textured.png
var tileGroundTextured []byte

//go:embed img/tile_ground_grass.png
var tileGroundGrass []byte

//go:embed img/tile_tree.png
var tileTree []byte

//go:embed img/tile_column.png
var tileColumn []byte

var (
	TileGroundTextured *ebiten.Image
	TileGroundGrass    *ebiten.Image
	TileTree           *ebiten.Image
	TileColumn         *ebiten.Image
	TileSize           = 16
)

func init() {
	loadTile := func(data []byte) *ebiten.Image {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println("Error loading tile:", err)
			return nil
		}
		return ebiten.NewImageFromImage(img)
	}

	TileGroundTextured = loadTile(tileGroundTextured)
	TileGroundGrass = loadTile(tileGroundGrass)
	TileTree = loadTile(tileTree)
	TileColumn = loadTile(tileColumn)
}
