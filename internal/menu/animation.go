package menu

import (
	"math"
	"math/rand"
)

const (
	initialObjectCount = 6
	objectTypeCount    = 5
	minVelocity        = 0.2
	velocityRange      = 0.3
	rotSpeedRange      = 0.03
	minAlpha           = 0.4
	alphaRange         = 0.2
	spawnChance        = 0.015
	spawnY             = -15.0
)

const (
	easterEggObjectCount = 25
	easterEggMinVel      = 0.8
	easterEggVelRange    = 1.5
	easterEggRotSpeed    = 0.08
	easterEggMinAlpha    = 0.6
	easterEggAlphaRange  = 0.3
)

const (
	subtitleText = "~ by VVV ~"
)

func (m *Menu) initSubtitleLetters() {
	m.subtitleText = subtitleText
	m.subtitleLetters = make([]SubtitleLetter, len(m.subtitleText))
	m.subtitleShowPhase = SubtitleAppearing

	for i := range m.subtitleLetters {
		m.subtitleLetters[i] = SubtitleLetter{
			OffsetX:   (rand.Float64() - 0.5) * 200,
			OffsetY:   (rand.Float64() - 0.5) * 100,
			Alpha:     0,
			Rotation:  (rand.Float64() - 0.5) * 2,
			Scale:     0.5 + rand.Float64()*1.5,
			Visible:   false,
			EchoAlpha: 0,
		}
	}
}

func (m *Menu) triggerSubtitleGlitch() {
	m.subtitleShowPhase = SubtitleGlitching
	m.subtitleAnimTimer = 0

	for i := range m.subtitleLetters {
		m.subtitleLetters[i].TargetX = (rand.Float64() - 0.5) * 150
		m.subtitleLetters[i].TargetY = (rand.Float64() - 0.5) * 80
		m.subtitleLetters[i].EchoAlpha = 0.8
	}
}

func (m *Menu) updateSubtitleLetters() {
	m.subtitleAnimTimer++

	switch m.subtitleShowPhase {
	case SubtitleAppearing:
		m.updateSubtitleAppearing()
	case SubtitleVisible:
		m.updateSubtitleVisible()
	case SubtitleGlitching:
		m.updateSubtitleGlitching()
	}
}

func (m *Menu) updateSubtitleAppearing() {
	allVisible := true
	for i := range m.subtitleLetters {
		letter := &m.subtitleLetters[i]

		// Staggered appearance
		appearDelay := i * LetterAppearDelay
		if m.subtitleAnimTimer > appearDelay {
			letter.Visible = true
			// Lerp to center
			letter.OffsetX *= 0.9
			letter.OffsetY *= 0.9
			letter.Rotation *= 0.9
			letter.Scale = letter.Scale*0.9 + 1.0*0.1
			letter.Alpha = math.Min(1.0, letter.Alpha+0.05)
		}

		if !letter.Visible || letter.Alpha < 0.99 {
			allVisible = false
		}
	}

	if allVisible && m.subtitleAnimTimer > len(m.subtitleLetters)*LetterAppearDelay+30 {
		m.subtitleShowPhase = SubtitleVisible
	}
}

func (m *Menu) updateSubtitleVisible() {
	for i := range m.subtitleLetters {
		letter := &m.subtitleLetters[i]
		wave := math.Sin(float64(m.frameCount)*0.05 + float64(i)*0.3)
		letter.OffsetY = wave * 2
		letter.EchoAlpha *= 0.95
	}
}

func (m *Menu) updateSubtitleGlitching() {
	const scatterDuration = 20

	if m.subtitleAnimTimer < scatterDuration {
		// Scatter phase
		for i := range m.subtitleLetters {
			letter := &m.subtitleLetters[i]
			letter.OffsetX = letter.OffsetX*0.8 + letter.TargetX*0.2
			letter.OffsetY = letter.OffsetY*0.8 + letter.TargetY*0.2
			letter.Rotation += (rand.Float64() - 0.5) * 0.3
			letter.Alpha = 0.5 + rand.Float64()*0.5
			letter.EchoAlpha = 0.6
		}
	} else {
		// Return phase
		allReturned := true
		for i := range m.subtitleLetters {
			letter := &m.subtitleLetters[i]
			letter.OffsetX *= 0.85
			letter.OffsetY *= 0.85
			letter.Rotation *= 0.85
			letter.Alpha = letter.Alpha*0.9 + 1.0*0.1
			letter.EchoAlpha *= 0.92

			if math.Abs(letter.OffsetX) > 1 || math.Abs(letter.OffsetY) > 1 {
				allReturned = false
			}
		}

		if allReturned && m.subtitleAnimTimer > LetterReturnDuration {
			m.subtitleShowPhase = SubtitleVisible
		}
	}
}

func (m *Menu) spawnInitialObjects() {
	for i := 0; i < initialObjectCount; i++ {
		objType := rand.Intn(objectTypeCount)
		m.objects = append(m.objects, FallingObject{
			X:        rand.Float64() * ScreenWidth,
			Y:        rand.Float64() * ScreenHeight,
			VelY:     minVelocity + rand.Float64()*velocityRange,
			ObjType:  objType,
			Rotation: rand.Float64() * math.Pi * 2,
			RotSpeed: (rand.Float64() - 0.5) * rotSpeedRange,
			Alpha:    minAlpha + rand.Float64()*alphaRange,
		})
	}
}

func (m *Menu) spawnObject() {
	if rand.Float64() < spawnChance {
		objType := rand.Intn(objectTypeCount)
		m.objects = append(m.objects, FallingObject{
			X:        rand.Float64() * ScreenWidth,
			Y:        spawnY,
			VelY:     minVelocity + rand.Float64()*velocityRange,
			ObjType:  objType,
			Rotation: rand.Float64() * math.Pi * 2,
			RotSpeed: (rand.Float64() - 0.5) * rotSpeedRange,
			Alpha:    minAlpha + rand.Float64()*alphaRange,
		})
	}
}

func (m *Menu) spawnEasterEggObjects() {
	for i := 0; i < easterEggObjectCount; i++ {
		m.objects = append(m.objects, FallingObject{
			X:        rand.Float64() * ScreenWidth,
			Y:        rand.Float64()*-100 - 10,
			VelY:     easterEggMinVel + rand.Float64()*easterEggVelRange,
			ObjType:  rand.Intn(4),
			Rotation: rand.Float64() * math.Pi * 2,
			RotSpeed: (rand.Float64() - 0.5) * easterEggRotSpeed,
			Alpha:    easterEggMinAlpha + rand.Float64()*easterEggAlphaRange,
		})
	}
}

func (m *Menu) updateFallingObjects() {
	objectCleanupY := float64(ScreenHeight) + 20.0
	newObjects := make([]FallingObject, 0, len(m.objects))
	for _, obj := range m.objects {
		obj.Y += obj.VelY
		obj.Rotation += obj.RotSpeed
		if obj.Y < objectCleanupY {
			newObjects = append(newObjects, obj)
		}
	}
	m.objects = newObjects
}

func (m *Menu) pulseValue() float64 {
	return 0.85 + math.Sin(float64(m.frameCount)*0.1)*0.15
}
