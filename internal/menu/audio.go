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
	shouldPlayMenu := m.shouldPlayMenuMusic() && !audio.IsMuted()
	shouldPlayVictory := m.state == StateLevelComplete && !audio.IsMuted()
	shouldPlayGameOver := m.state == StateGameOver && !audio.IsMuted()

	if shouldPlayMenu && !audio.IsMusicPlaying(audio.SoundMenuMusic) {
		audio.PlayMusic(audio.SoundMenuMusic)
	} else if !shouldPlayMenu && audio.IsMusicPlaying(audio.SoundMenuMusic) {
		audio.StopMusic(audio.SoundMenuMusic)
	}

	// HELL YEAAAH BROTHER AND SISTERS HELL YEAAAH
	if shouldPlayVictory && !audio.IsMusicPlaying(audio.SoundVictory) {
		audio.PlayMusic(audio.SoundVictory)
	} else if !shouldPlayVictory && audio.IsMusicPlaying(audio.SoundVictory) {
		audio.StopMusic(audio.SoundVictory)
	}

	if shouldPlayGameOver && !audio.IsMusicPlaying(audio.SoundGameOver) {
		audio.PlayMusic(audio.SoundGameOver)
	} else if !shouldPlayGameOver && audio.IsMusicPlaying(audio.SoundGameOver) {
		audio.StopMusic(audio.SoundGameOver)
	}

	audio.UpdateMusicVolume()
}

func (m *Menu) shouldPlayMenuMusic() bool {
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
