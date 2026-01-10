package menu

import (
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
)

func (m *Menu) playSelectSound() {
	audio.Play(audio.SoundMenuSelect)
}

func (m *Menu) playConfirmSound() {
	audio.Play(audio.SoundMenuConfirm)
}

func (m *Menu) startMusic() {
	if m.musicPlaying {
		return
	}
	audio.PlayMusic(audio.SoundMenuMusic)
	m.musicPlaying = true
}

func (m *Menu) stopMusic() {
	audio.StopMusic(audio.SoundMenuMusic)
	m.musicPlaying = false
}

func (m *Menu) updateMusicState() {
	shouldPlay := m.shouldPlayMusic() && !audio.IsMuted()

	if shouldPlay && !audio.IsMusicPlaying(audio.SoundMenuMusic) {
		audio.PlayMusic(audio.SoundMenuMusic)
	} else if !shouldPlay && audio.IsMusicPlaying(audio.SoundMenuMusic) {
		audio.StopMusic(audio.SoundMenuMusic)
	}

	audio.UpdateMusicVolume()
}

func (m *Menu) shouldPlayMusic() bool {
	switch m.state {
	case StateMenu:
		return true
	case StateSettings:
		return m.previousState == StateMenu
	default:
		return false
	}
}

func (m *Menu) updateMusicVolume() {
	audio.UpdateMusicVolume()
}
