package menu

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/fonts/PressStart2P-Regular.ttf
var fontData []byte

func FontData() []byte {
	return fontData
}

//go:embed assets/img/skull.png
var skullPNG []byte

//go:embed assets/img/rotten_apple.png
var rottenApplePNG []byte

//go:embed assets/img/dead_orange.png
var deadOrangePNG []byte

//go:embed assets/img/withered_cherry.png
var witheredCherryPNG []byte

//go:embed assets/img/rotted_banana.png
var rottedBananaPNG []byte

//go:embed assets/img/title_dead.png
var titleDeadPNG []byte

//go:embed assets/img/title_jump.png
var titleJumpPNG []byte

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func (m *Menu) loadAssets() {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(err)
	}
	m.fontFace = source
	m.fontSmall = &text.GoTextFace{Source: source, Size: FontSizeSmall}
	m.fontMedium = &text.GoTextFace{Source: source, Size: FontSizeMedium}

	m.titleDeadImg = loadImage(titleDeadPNG)
	m.titleJumpImg = loadImage(titleJumpPNG)
	m.skullImg = loadImage(skullPNG)

	m.objectImages = []*ebiten.Image{
		loadImage(rottenApplePNG),
		loadImage(deadOrangePNG),
		loadImage(witheredCherryPNG),
		loadImage(rottedBananaPNG),
	}
}
