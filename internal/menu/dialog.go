package menu

import (
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (m *Menu) ShowLevelComplete() {
	m.state = StateLevelComplete
	m.selectedIndex = 0
	audio.RestartMusic(audio.SoundVictory)
	if m.OnLevelComplete != nil {
		m.OnLevelComplete()
	}
}

func (m *Menu) ShowGameOver() {
	m.state = StateGameOver
	m.selectedIndex = 0
	audio.RestartMusic(audio.SoundGameOver)
	if m.OnGameOver != nil {
		m.OnGameOver()
	}
}

func (m *Menu) ShowEpilogueEnding() {
	m.state = StateEpilogueEnding
	m.selectedIndex = 0
	m.epilogueTimer = 0
}

func (m *Menu) ShowDialog(title, message string, buttons []MenuItem) {
	m.activeDialog = &Dialog{
		Title:    title,
		Message:  message,
		Buttons:  buttons,
		Selected: 0,
	}
}

func (m *Menu) CloseDialog() {
	m.activeDialog = nil
}

func (m *Menu) HasActiveDialog() bool {
	return m.activeDialog != nil
}

func (m *Menu) updateDialog() {
	if m.activeDialog == nil {
		return
	}

	if inpututil.IsKeyJustPressed(keyLeft) || inpututil.IsKeyJustPressed(keyA) {
		m.activeDialog.Selected--
		if m.activeDialog.Selected < 0 {
			m.activeDialog.Selected = len(m.activeDialog.Buttons) - 1
		}
		m.playSelectSound()
	}
	if inpututil.IsKeyJustPressed(keyRight) || inpututil.IsKeyJustPressed(keyD) {
		m.activeDialog.Selected++
		if m.activeDialog.Selected >= len(m.activeDialog.Buttons) {
			m.activeDialog.Selected = 0
		}
		m.playSelectSound()
	}

	if inpututil.IsKeyJustPressed(keyEnter) || inpututil.IsKeyJustPressed(keySpace) {
		m.playConfirmSound()
		if m.activeDialog.Buttons[m.activeDialog.Selected].Action != nil {
			m.activeDialog.Buttons[m.activeDialog.Selected].Action()
		}
	}

	if inpututil.IsKeyJustPressed(keyEscape) {
		m.playSelectSound()
		m.CloseDialog()
	}
}
