package systems

import (
	"github.com/game-jam-2026/dead-jump/internal/menu"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
)

var prevMusicState menu.GameState

func UpdateLevelMusic(currentState menu.GameState) {
	if currentState == prevMusicState {
		return
	}

	if currentState == menu.StatePlaying {
		audio.PlayMusic(audio.SoundLevelMusic)
	}

	if prevMusicState == menu.StatePlaying && currentState != menu.StatePlaying {
		audio.StopMusic(audio.SoundLevelMusic)
	}

	prevMusicState = currentState
}
