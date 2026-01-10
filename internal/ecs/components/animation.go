package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Duration            time.Duration
	Images              []*ebiten.Image
	lastChangeTimestamp time.Time
	lastImage           int
}

func (a *Animation) CheckAndGetImage() *ebiten.Image {
	if time.Now().Compare(a.lastChangeTimestamp.Add(a.Duration)) < 1 {
		return nil
	}

	var img *ebiten.Image
	if a.lastImage+1 < len(a.Images) {
		img = a.Images[a.lastImage+1]
		a.lastImage = a.lastImage + 1
	} else {
		img = a.Images[0]
		a.lastImage = 0
	}

	a.lastChangeTimestamp = time.Now()

	return img
}
