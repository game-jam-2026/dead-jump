package menu

import (
	"math"
	"math/rand"

	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	keyUp     = ebiten.KeyUp
	keyDown   = ebiten.KeyDown
	keyLeft   = ebiten.KeyLeft
	keyRight  = ebiten.KeyRight
	keyW      = ebiten.KeyW
	keyS      = ebiten.KeyS
	keyA      = ebiten.KeyA
	keyD      = ebiten.KeyD
	keyEnter  = ebiten.KeyEnter
	keySpace  = ebiten.KeySpace
	keyEscape = ebiten.KeyEscape
)

func NewMenu() *Menu {
	m := &Menu{
		state:         StateMenu,
		selectedIndex: 0,
		objects:       make([]FallingObject, 0),
		secretCode:    []ebiten.Key{ebiten.KeyF, ebiten.KeyR, ebiten.KeyU, ebiten.KeyI, ebiten.KeyT},
		secretIndex:   0,
	}

	m.initMenuItems()
	m.initPauseItems()
	m.initConfirmItems()
	m.initSettingsItems()
	m.initLevelCompleteItems()
	m.initGameOverItems()
	m.initEpilogueItems()

	m.loadAssets()
	m.spawnInitialObjects()
	m.initSubtitleLetters()

	return m
}

func (m *Menu) initMenuItems() {
	m.items = []MenuItem{
		{Text: "START", Action: func() {
			if m.OnStartGame != nil {
				m.OnStartGame()
			}
		}},
		{Text: "SETTINGS", Action: func() {
			m.previousState = StateMenu
			m.state = StateSettings
			m.selectedIndex = 0
		}},
		{Text: "QUIT", Action: func() {
			if m.OnQuit != nil {
				m.OnQuit()
			}
		}},
	}
}

func (m *Menu) initPauseItems() {
	m.pauseItems = []MenuItem{
		{Text: "RESUME", Action: func() {
			if m.OnResume != nil {
				m.OnResume()
			}
		}},
		{Text: "RESTART", Action: func() {
			m.state = StateConfirmRestart
			m.selectedIndex = 1
		}},
		{Text: "SETTINGS", Action: func() {
			m.previousState = StatePaused
			m.state = StateSettings
			m.selectedIndex = 0
		}},
		{Text: "MAIN MENU", Action: func() {
			m.state = StateMenu
			m.selectedIndex = 0
		}},
		{Text: "QUIT", Action: func() {
			if m.OnQuit != nil {
				m.OnQuit()
			}
		}},
	}
}

func (m *Menu) initConfirmItems() {
	m.confirmItems = []MenuItem{
		{Text: "YES", Action: func() {
			if m.OnRestart != nil {
				m.OnRestart()
			}
			m.SetState(StatePlaying)
		}},
		{Text: "NO", Action: func() {
			m.state = StatePaused
			m.selectedIndex = 1
		}},
	}
}

func (m *Menu) initLevelCompleteItems() {
	m.levelCompleteItems = []MenuItem{
		{Text: "NEXT LEVEL", Action: func() {
			audio.StopMusic(audio.SoundVictory)
			if m.OnNextLevel != nil {
				m.OnNextLevel()
			}
			m.SetState(StatePlaying)
		}},
		{Text: "RESTART", Action: func() {
			audio.StopMusic(audio.SoundVictory)
			if m.OnRestart != nil {
				m.OnRestart()
			}
			m.SetState(StatePlaying)
		}},
		{Text: "MAIN MENU", Action: func() {
			audio.StopMusic(audio.SoundVictory)
			if m.OnMainMenu != nil {
				m.OnMainMenu()
			}
			m.state = StateMenu
			m.selectedIndex = 0
		}},
	}
}

func (m *Menu) initGameOverItems() {
	m.gameOverItems = []MenuItem{
		{Text: "TRY AGAIN", Action: func() {
			audio.StopMusic(audio.SoundGameOver)
			if m.OnRestart != nil {
				m.OnRestart()
			}
			m.SetState(StatePlaying)
		}},
		{Text: "MAIN MENU", Action: func() {
			audio.StopMusic(audio.SoundGameOver)
			if m.OnMainMenu != nil {
				m.OnMainMenu()
			}
			m.state = StateMenu
			m.selectedIndex = 0
		}},
	}
}

func (m *Menu) initEpilogueItems() {
	m.epilogueItems = []MenuItem{
		{Text: "CONTINUE", Action: func() {
			if m.OnEpilogueComplete != nil {
				m.OnEpilogueComplete()
			}
			m.SetState(StatePlaying)
		}},
	}
}

func (m *Menu) initSettingsItems() {
	m.settingsItems = []MenuItem{
		{Text: m.getMasterVolumeText(), Action: func() {}},
		{Text: m.getMusicVolumeText(), Action: func() {}},
		{Text: m.getSFXVolumeText(), Action: func() {}},
		{Text: "BACK", Action: func() {
			m.state = m.previousState
			m.selectedIndex = 0
		}},
	}
}

func (m *Menu) getMasterVolumeText() string {
	vol := int(audio.GetMasterVolume() * VolumeSteps)
	return "MASTER:  " + m.volumeBar(vol)
}

func (m *Menu) getMusicVolumeText() string {
	vol := int(audio.GetMusicVolume() * VolumeSteps)
	return "MUSIC:   " + m.volumeBar(vol)
}

func (m *Menu) getSFXVolumeText() string {
	vol := int(audio.GetSFXVolume() * VolumeSteps)
	return "SFX:     " + m.volumeBar(vol)
}

func (m *Menu) volumeBar(level int) string {
	bar := "<"
	for i := 0; i < VolumeSteps; i++ {
		if i < level {
			bar += "#"
		} else {
			bar += "-"
		}
	}
	bar += ">"
	return bar
}

func (m *Menu) updateSettingsItems() {
	if len(m.settingsItems) >= 3 {
		m.settingsItems[0].Text = m.getMasterVolumeText()
		m.settingsItems[1].Text = m.getMusicVolumeText()
		m.settingsItems[2].Text = m.getSFXVolumeText()
	}
}

func (m *Menu) adjustVolume(delta float64) {
	switch m.selectedIndex {
	case 0:
		audio.SetMasterVolume(audio.GetMasterVolume() + delta)
	case 1:
		audio.SetMusicVolume(audio.GetMusicVolume() + delta)
	case 2:
		audio.SetSFXVolume(audio.GetSFXVolume() + delta)
	}
	m.updateSettingsItems()
	m.updateMusicVolume()
}

func (m *Menu) GetState() GameState {
	return m.state
}

func (m *Menu) SetState(state GameState) {
	m.state = state
	m.selectedIndex = 0

	if state == StatePlaying {
		m.stopMusic()
	}
}

func (m *Menu) Update() {
	m.frameCount++
	m.titleOffset = math.Sin(float64(m.frameCount)*0.03) * 2

	m.updateMusicState()
	m.updateScreenShake()
	m.updateGlitchEffect()
	m.updateSubtitleLetters()
	m.updateFallingObjectsWithSpawn()
	m.handleEasterEgg()

	if m.activeDialog != nil {
		m.updateDialog()
		return
	}

	if m.state == StateEpilogueEnding {
		m.updateEpilogue()
		return
	}

	items := m.getActiveItems()
	if items == nil {
		return
	}

	m.handleNavigation(items)
	m.handleVolumeAdjustment()
	m.handleConfirmAction(items)
	m.handleEscapeKey()
}

func (m *Menu) updateEpilogue() {
	m.epilogueTimer++

	if m.epilogueTimer > 60 {
		if inpututil.IsKeyJustPressed(keyEnter) || inpututil.IsKeyJustPressed(keySpace) {
			m.playConfirmSound()
			if len(m.epilogueItems) > 0 && m.epilogueItems[0].Action != nil {
				m.epilogueItems[0].Action()
			}
		}
	}
}

func (m *Menu) updateScreenShake() {
	if m.screenShake > 0 {
		m.screenShake *= 0.9
	}
}

func (m *Menu) updateGlitchEffect() {
	m.dejavuTimer++
	if m.dejavuTimer > GlitchMinInterval && rand.Float64() < GlitchChance {
		m.dejavuGlitch = 8
		m.dejavuTimer = 0
		m.screenShake = 3
		m.triggerSubtitleGlitch()
	}
	if m.dejavuGlitch > 0 {
		m.dejavuGlitch -= 0.4
	}
}

func (m *Menu) updateFallingObjectsWithSpawn() {
	m.spawnObject()
	m.updateFallingObjects()
}

func (m *Menu) handleEasterEgg() {
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		if key == m.secretCode[m.secretIndex] {
			m.secretIndex++
			if m.secretIndex >= len(m.secretCode) {
				m.activateEasterEgg()
			}
		} else {
			m.secretIndex = 0
		}
	}

	if m.easterEggTimer > 0 {
		m.easterEggTimer--
		if m.easterEggTimer == 0 {
			m.easterEggActive = false
		}
	}
}

func (m *Menu) activateEasterEgg() {
	m.easterEggActive = true
	m.easterEggTimer = EasterEggDuration
	m.secretIndex = 0
	m.screenShake = 10
	m.spawnEasterEggObjects()
}

func (m *Menu) getActiveItems() []MenuItem {
	switch m.state {
	case StateMenu:
		return m.items
	case StatePaused:
		return m.pauseItems
	case StateConfirmRestart:
		return m.confirmItems
	case StateSettings:
		return m.settingsItems
	case StateLevelComplete:
		return m.levelCompleteItems
	case StateGameOver:
		return m.gameOverItems
	case StateEpilogueEnding:
		return m.epilogueItems
	case StatePlaying:
		return nil
	default:
		return nil
	}
}

func (m *Menu) handleNavigation(items []MenuItem) {
	if inpututil.IsKeyJustPressed(keyUp) || inpututil.IsKeyJustPressed(keyW) {
		m.selectedIndex--
		if m.selectedIndex < 0 {
			m.selectedIndex = len(items) - 1
		}
		m.playSelectSound()
	}
	if inpututil.IsKeyJustPressed(keyDown) || inpututil.IsKeyJustPressed(keyS) {
		m.selectedIndex++
		if m.selectedIndex >= len(items) {
			m.selectedIndex = 0
		}
		m.playSelectSound()
	}
}

func (m *Menu) handleVolumeAdjustment() {
	if m.state != StateSettings || m.selectedIndex >= 3 {
		return
	}

	if inpututil.IsKeyJustPressed(keyLeft) || inpututil.IsKeyJustPressed(keyA) {
		m.adjustVolume(-VolumeStepValue)
		m.playSelectSound()
	}
	if inpututil.IsKeyJustPressed(keyRight) || inpututil.IsKeyJustPressed(keyD) {
		m.adjustVolume(VolumeStepValue)
		m.playSelectSound()
	}
}

func (m *Menu) handleConfirmAction(items []MenuItem) {
	if inpututil.IsKeyJustPressed(keyEnter) || inpututil.IsKeyJustPressed(keySpace) {
		m.playConfirmSound()
		if items[m.selectedIndex].Action != nil {
			items[m.selectedIndex].Action()
		}
	}
}

func (m *Menu) handleEscapeKey() {
	if !inpututil.IsKeyJustPressed(keyEscape) {
		return
	}

	switch m.state {
	case StateConfirmRestart:
		m.playSelectSound()
		m.state = StatePaused
		m.selectedIndex = 1
	case StateSettings:
		m.playSelectSound()
		m.state = m.previousState
		m.selectedIndex = 0
	}
}
