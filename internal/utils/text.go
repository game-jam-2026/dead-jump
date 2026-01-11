package utils

import (
	"bytes"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/menu"
)

var loreFont *text.GoTextFace

func init() {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(menu.FontData()))
	if err != nil {
		panic(err)
	}
	loreFont = &text.GoTextFace{Source: fontSource, Size: 8}
}

func DrawLoreText(w *ecs.World, screen *ebiten.Image) {
	loreText, err := ecs.GetResource[components.LoreText](w)
	if err != nil || loreText.Text == "" {
		return
	}

	maxWidth := 280
	lines := WrapText(loreText.Text, maxWidth, loreFont)

	startY := 30.0
	lineHeight := 14.0
	centerX := float64(menu.ScreenWidth) / 2

	for i, line := range lines {
		op := &text.DrawOptions{}
		op.GeoM.Translate(centerX, startY+float64(i)*lineHeight)
		op.PrimaryAlign = text.AlignCenter
		op.ColorScale.ScaleWithColor(color.RGBA{200, 180, 160, 255})
		text.Draw(screen, line, loreFont, op)
	}
}

func WrapText(txt string, maxWidth int, face *text.GoTextFace) []string {
	words := strings.Fields(txt)
	var lines []string
	var currentLine string

	charWidth := face.Size
	maxChars := int(float64(maxWidth) / charWidth)

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) > maxChars && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
