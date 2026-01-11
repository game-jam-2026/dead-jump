package menu

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (m *Menu) Draw(screen *ebiten.Image) {
	shakeX, shakeY := m.getScreenShake()

	if m.state == StateEpilogueEnding {
		m.drawEpilogueEndingScreen(screen, shakeX, shakeY)
		return
	}

	isInGameOverlay := m.state == StateLevelComplete || m.state == StateGameOver ||
		m.state == StatePaused || m.state == StateConfirmRestart ||
		(m.state == StateSettings && m.previousState == StatePaused)

	if isInGameOverlay {
		m.drawDarkOverlay(screen)
		m.drawStateContent(screen, shakeX, shakeY)
		if m.activeDialog != nil {
			m.drawActiveDialog(screen, shakeX, shakeY)
		}
		return
	}

	screen.Fill(colorBgDark)

	m.drawFallingObjects(screen, shakeX, shakeY)
	m.drawTitle(screen, shakeX, shakeY)
	m.drawSubtitle(screen, shakeX, shakeY)
	m.drawStateContent(screen, shakeX, shakeY)

	if m.activeDialog != nil {
		m.drawActiveDialog(screen, shakeX, shakeY)
	}

	m.drawHint(screen, shakeY)

	if m.easterEggActive {
		centerX := float64(ScreenWidth) / 2
		m.drawText(screen, "ROTTEN RAIN!", centerX, 185+shakeY, m.fontSmall, colorDeathRed, true)
	}
}

func (m *Menu) getScreenShake() (float64, float64) {
	if m.screenShake > 0.5 {
		return (rand.Float64() - 0.5) * m.screenShake, (rand.Float64() - 0.5) * m.screenShake
	}
	return 0, 0
}

func (m *Menu) drawFallingObjects(screen *ebiten.Image, shakeX, shakeY float64) {
	for _, obj := range m.objects {
		var img *ebiten.Image
		if obj.ObjType < 4 && obj.ObjType < len(m.objectImages) {
			img = m.objectImages[obj.ObjType]
		} else {
			img = m.skullImg
		}
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			w, h := img.Bounds().Dx(), img.Bounds().Dy()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(obj.Rotation)
			op.GeoM.Translate(obj.X+shakeX, obj.Y+shakeY)
			op.ColorScale.ScaleAlpha(float32(obj.Alpha))
			screen.DrawImage(img, op)
		}
	}
}

func (m *Menu) drawStateContent(screen *ebiten.Image, shakeX, shakeY float64) {
	switch m.state {
	case StateMenu:
		m.drawMenuItems(screen, m.items, 120, shakeX, shakeY)
	case StatePaused:
		m.drawPauseOverlay(screen)
		// Menu starts 20px after title, title is at (ScreenHeight-120)/2 = 60
		menuStartY := (float64(ScreenHeight)-120)/2 + 20
		m.drawMenuItems(screen, m.pauseItems, menuStartY, shakeX, shakeY)
	case StateConfirmRestart:
		m.drawPauseOverlay(screen)
		m.drawConfirmDialog(screen, shakeX, shakeY)
	case StateSettings:
		m.drawDarkOverlay(screen)
		m.drawSettingsMenu(screen, shakeX, shakeY)
	case StateLevelComplete:
		m.drawLevelCompleteScreen(screen, shakeX, shakeY)
	case StateGameOver:
		m.drawGameOverScreen(screen, shakeX, shakeY)
	case StatePlaying:
	}
}

func (m *Menu) drawHint(screen *ebiten.Image, shakeY float64) {
	centerX := float64(ScreenWidth) / 2
	hintY := float64(ScreenHeight) - 22 + shakeY
	var hint string
	if m.state == StateSettings && m.selectedIndex < 3 {
		hint = "< > VOLUME  ESC BACK"
	} else {
		hint = "ARROWS + ENTER"
	}
	hintColor := color.RGBA{R: 50, G: 45, B: 60, A: 255}
	m.drawText(screen, hint, centerX, hintY, m.fontSmall, hintColor, true)
}

func (m *Menu) drawTitle(screen *ebiten.Image, shakeX, shakeY float64) {
	deadW := float64(m.titleDeadImg.Bounds().Dx())
	jumpW := float64(m.titleJumpImg.Bounds().Dx())
	totalW := deadW + TitleSpacing + jumpW
	startX := (ScreenWidth - totalW) / 2
	titleY := 30 + m.titleOffset + shakeY

	const glitchOffset = 3.0
	if m.dejavuGlitch > 0 || m.frameCount%100 < 8 {
		m.drawTitleGlitch(screen, startX, deadW, titleY, shakeX, glitchOffset)
	}

	opDead := &ebiten.DrawImageOptions{}
	opDead.GeoM.Translate(startX+shakeX, titleY)
	screen.DrawImage(m.titleDeadImg, opDead)

	opJump := &ebiten.DrawImageOptions{}
	opJump.GeoM.Translate(startX+deadW+TitleSpacing+shakeX, titleY)
	screen.DrawImage(m.titleJumpImg, opJump)
}

func (m *Menu) drawTitleGlitch(screen *ebiten.Image, startX, deadW, titleY, shakeX, glitchOffset float64) {
	// Red ghost
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(startX-m.dejavuGlitch-glitchOffset+shakeX, titleY+2)
	op.ColorScale.Scale(1.5, 0.3, 0.3, 0.4)
	screen.DrawImage(m.titleDeadImg, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(startX+deadW+TitleSpacing-m.dejavuGlitch-glitchOffset+shakeX, titleY+2)
	op2.ColorScale.Scale(1.5, 0.3, 0.3, 0.4)
	screen.DrawImage(m.titleJumpImg, op2)

	// Blue ghost
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(startX+m.dejavuGlitch+glitchOffset+shakeX, titleY-2)
	op3.ColorScale.Scale(0.3, 0.3, 1.5, 0.4)
	screen.DrawImage(m.titleDeadImg, op3)

	op4 := &ebiten.DrawImageOptions{}
	op4.GeoM.Translate(startX+deadW+TitleSpacing+m.dejavuGlitch+glitchOffset+shakeX, titleY-2)
	op4.ColorScale.Scale(0.3, 0.3, 1.5, 0.4)
	screen.DrawImage(m.titleJumpImg, op4)
}

func (m *Menu) drawSubtitle(screen *ebiten.Image, shakeX, shakeY float64) {
	baseX := float64(ScreenWidth) / 2
	baseY := 70.0
	charWidth := CharWidthPixels

	totalWidth := float64(len(m.subtitleText)) * charWidth
	startX := baseX - totalWidth/2

	for i, ch := range m.subtitleText {
		if i >= len(m.subtitleLetters) {
			break
		}
		letter := &m.subtitleLetters[i]

		if !letter.Visible && m.subtitleShowPhase != SubtitleAppearing {
			continue
		}

		charX := startX + float64(i)*charWidth + letter.OffsetX + shakeX
		charY := baseY + letter.OffsetY + shakeY

		// Draw echo effect
		if letter.EchoAlpha > 0.05 {
			m.drawLetterEcho(screen, ch, letter, startX, baseY, float64(i), charWidth, shakeX, shakeY)
		}

		// Draw main letter
		m.drawLetter(screen, ch, letter, charX, charY)
	}
}

func (m *Menu) drawLetterEcho(screen *ebiten.Image, ch rune, letter *SubtitleLetter, startX, baseY, i, charWidth, shakeX, shakeY float64) {
	for e := 3; e >= 1; e-- {
		echoOp := &text.DrawOptions{}
		echoOffsetX := letter.OffsetX * float64(e) * 0.3
		echoOffsetY := letter.OffsetY * float64(e) * 0.3
		echoOp.GeoM.Translate(startX+i*charWidth+echoOffsetX+shakeX, baseY+echoOffsetY+shakeY)

		alpha := letter.EchoAlpha * (0.3 / float64(e))
		echoColor := color.RGBA{
			R: uint8(float64(colorDeadPurple.R) * 0.7),
			G: uint8(float64(colorDeadPurple.G) * 0.5),
			B: uint8(float64(colorDeadPurple.B) * 1.2),
			A: uint8(255 * alpha),
		}
		echoOp.ColorScale.ScaleWithColor(echoColor)
		text.Draw(screen, string(ch), m.fontSmall, echoOp)
	}
}

func (m *Menu) drawLetter(screen *ebiten.Image, ch rune, letter *SubtitleLetter, charX, charY float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(charX, charY)

	alpha := letter.Alpha
	var c color.RGBA
	if m.subtitleShowPhase == SubtitleGlitching {
		// Glitching - flicker between colors
		if rand.Float64() < 0.3 {
			c = color.RGBA{R: 180, G: 60, B: 80, A: uint8(255 * alpha)}
		} else {
			c = color.RGBA{R: 80, G: 60, B: 180, A: uint8(255 * alpha)}
		}
	} else {
		c = color.RGBA{
			R: uint8(float64(colorDeadPurple.R) * alpha),
			G: uint8(float64(colorDeadPurple.G) * alpha),
			B: uint8(float64(colorDeadPurple.B) * alpha),
			A: uint8(255 * alpha),
		}
	}

	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, string(ch), m.fontSmall, op)
}

func (m *Menu) drawText(screen *ebiten.Image, txt string, x, y float64, face *text.GoTextFace, c color.Color, centered bool) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	if centered {
		op.PrimaryAlign = text.AlignCenter
	}
	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, txt, face, op)
}

func (m *Menu) drawMenuItems(screen *ebiten.Image, items []MenuItem, startY float64, shakeX, shakeY float64) {
	centerX := float64(ScreenWidth) / 2

	for i, item := range items {
		y := startY + float64(i)*MenuItemSpacing + shakeY
		selected := i == m.selectedIndex

		c := colorDimGray
		if selected {
			c = m.pulseColor(colorSelectedGlow)
			m.drawSkullIndicators(screen, item.Text, centerX, y, shakeX)
		}

		m.drawText(screen, item.Text, centerX+shakeX, y, m.fontSmall, c, true)
	}
}

func (m *Menu) drawSkullIndicators(screen *ebiten.Image, itemText string, centerX, y, shakeX float64) {
	if m.skullImg == nil {
		return
	}
	skullPulse := math.Sin(float64(m.frameCount)*0.12) * 3
	textW := float64(len(itemText)) * CharWidthPixels

	leftOp := &ebiten.DrawImageOptions{}
	leftOp.GeoM.Translate(centerX-textW/2-SkullMargin+skullPulse+shakeX, y-3)
	leftOp.ColorScale.ScaleAlpha(0.8)
	screen.DrawImage(m.skullImg, leftOp)

	rightOp := &ebiten.DrawImageOptions{}
	rightOp.GeoM.Scale(-1, 1)
	rightOp.GeoM.Translate(centerX+textW/2+SkullMargin-skullPulse+shakeX, y-3)
	rightOp.ColorScale.ScaleAlpha(0.8)
	screen.DrawImage(m.skullImg, rightOp)
}

func (m *Menu) drawDarkOverlay(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{R: 5, G: 3, B: 10, A: 180})
	screen.DrawImage(overlay, nil)
}

func (m *Menu) drawPauseOverlay(screen *ebiten.Image) {
	centerX := float64(ScreenWidth) / 2
	titleY := (float64(ScreenHeight) - 120) / 2
	m.drawText(screen, "PAUSED", centerX, titleY, m.fontMedium, colorBloodRed, true)
}

func (m *Menu) drawConfirmDialog(screen *ebiten.Image, shakeX, shakeY float64) {
	centerX := float64(ScreenWidth) / 2
	boxW := 180
	boxH := 75
	boxX := (ScreenWidth - boxW) / 2
	boxY := 85

	m.drawDialogBox(screen, boxX, boxY, boxW, boxH)

	m.drawText(screen, "RESTART?", centerX+shakeX, float64(boxY)+15+shakeY, m.fontSmall, colorWarning, true)
	m.drawText(screen, "PROGRESS LOST", centerX+shakeX, float64(boxY)+32+shakeY, m.fontSmall, colorDimGray, true)

	// YES / NO buttons
	buttonSpacing := 60.0
	startBtnX := centerX - buttonSpacing/2
	for i, item := range m.confirmItems {
		btnX := startBtnX + float64(i)*buttonSpacing + shakeX
		btnY := float64(boxY) + 52 + shakeY
		selected := i == m.selectedIndex

		c := colorDimGray
		if selected {
			if item.Text == "YES" {
				c = colorDeathRed
			} else {
				c = colorGhostWhite
			}
			m.drawText(screen, ">", btnX-12, btnY, m.fontSmall, c, false)
		}
		m.drawText(screen, item.Text, btnX, btnY, m.fontSmall, c, false)
	}
}

func (m *Menu) drawDialogBox(screen *ebiten.Image, boxX, boxY, boxW, boxH int) {
	borderColor := color.RGBA{R: 60, G: 30, B: 40, A: 255}
	fillColor := color.RGBA{R: 20, G: 15, B: 25, A: 255}
	borderWidth := 2

	// Border
	for x := boxX - borderWidth; x < boxX+boxW+borderWidth; x++ {
		for y := boxY - borderWidth; y < boxY+boxH+borderWidth; y++ {
			screen.Set(x, y, borderColor)
		}
	}
	// Fill
	for x := boxX; x < boxX+boxW; x++ {
		for y := boxY; y < boxY+boxH; y++ {
			screen.Set(x, y, fillColor)
		}
	}
}

func (m *Menu) drawSettingsMenu(screen *ebiten.Image, shakeX, shakeY float64) {
	centerX := float64(ScreenWidth) / 2
	volumeArrowLeft := 45.0
	volumeArrowRight := float64(ScreenWidth) - 50

	m.drawText(screen, "SETTINGS", centerX+shakeX, 85+shakeY, m.fontMedium, colorBloodRed, true)

	startY := 110.0
	volumeItemCount := 3
	for i, item := range m.settingsItems {
		y := startY + float64(i)*SettingsItemSpacing + shakeY
		selected := i == m.selectedIndex

		c := colorDimGray
		if selected {
			c = m.pulseColor(colorSelectedGlow)

			if i < volumeItemCount {
				m.drawText(screen, "<", volumeArrowLeft+shakeX, y, m.fontSmall, c, false)
				m.drawText(screen, ">", volumeArrowRight+shakeX, y, m.fontSmall, c, false)
			}
		}

		m.drawText(screen, item.Text, centerX+shakeX, y, m.fontSmall, c, true)
	}
}

func (m *Menu) drawLevelCompleteScreen(screen *ebiten.Image, shakeX, shakeY float64) {
	centerX := float64(ScreenWidth) / 2
	m.drawText(screen, "LEVEL COMPLETE!", centerX+shakeX, 70+shakeY, m.fontMedium, colorSelectedGlow, true)
	m.drawMenuItems(screen, m.levelCompleteItems, 110, shakeX, shakeY)
}

func (m *Menu) drawGameOverScreen(screen *ebiten.Image, shakeX, shakeY float64) {
	centerX := float64(ScreenWidth) / 2
	m.drawText(screen, "GAME OVER", centerX+shakeX, 70+shakeY, m.fontMedium, colorDeathRed, true)
	m.drawText(screen, "YOU DIED", centerX+shakeX, 90+shakeY, m.fontSmall, colorDimGray, true)
	m.drawMenuItems(screen, m.gameOverItems, 120, shakeX, shakeY)
}

func (m *Menu) drawEpilogueEndingScreen(screen *ebiten.Image, shakeX, shakeY float64) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	centerX := float64(ScreenWidth) / 2
	centerY := float64(ScreenHeight) / 2

	m.drawText(screen, "This should be stopped.", centerX+shakeX, centerY-20+shakeY, m.fontMedium, colorGhostWhite, true)
	m.drawText(screen, "...but the cycle continues.", centerX+shakeX, centerY+10+shakeY, m.fontSmall, colorDimGray, true)

	if m.epilogueTimer > 60 {
		m.drawText(screen, "Press ENTER to continue", centerX+shakeX, centerY+60+shakeY, m.fontSmall, colorDeadPurple, true)
	}
}

func (m *Menu) pulseColor(base color.RGBA) color.RGBA {
	pulse := m.pulseValue()
	return color.RGBA{
		R: uint8(float64(base.R) * pulse),
		G: uint8(float64(base.G) * pulse),
		B: uint8(float64(base.B) * pulse),
		A: base.A,
	}
}

func (m *Menu) drawActiveDialog(screen *ebiten.Image, shakeX, shakeY float64) {
	if m.activeDialog == nil {
		return
	}

	m.drawDarkOverlay(screen)

	centerX := float64(ScreenWidth) / 2
	boxW := 200
	boxH := 80
	boxX := (ScreenWidth - boxW) / 2
	boxY := (ScreenHeight - boxH) / 2

	m.drawDialogBox(screen, boxX, boxY, boxW, boxH)

	// Draw title
	m.drawText(screen, m.activeDialog.Title, centerX+shakeX, float64(boxY)+15+shakeY, m.fontSmall, colorWarning, true)

	// Draw message
	if m.activeDialog.Message != "" {
		m.drawText(screen, m.activeDialog.Message, centerX+shakeX, float64(boxY)+32+shakeY, m.fontSmall, colorDimGray, true)
	}

	// Draw buttons
	buttonCount := len(m.activeDialog.Buttons)
	if buttonCount > 0 {
		buttonSpacing := 60.0
		totalWidth := float64(buttonCount-1) * buttonSpacing
		startBtnX := centerX - totalWidth/2

		for i, btn := range m.activeDialog.Buttons {
			btnX := startBtnX + float64(i)*buttonSpacing + shakeX
			btnY := float64(boxY) + 55 + shakeY
			selected := i == m.activeDialog.Selected

			c := colorDimGray
			if selected {
				c = m.pulseColor(colorSelectedGlow)
				m.drawText(screen, ">", btnX-12, btnY, m.fontSmall, c, false)
			}
			m.drawText(screen, btn.Text, btnX, btnY, m.fontSmall, c, false)
		}
	}
}
