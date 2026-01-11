package assets

import (
	_ "embed"

	"github.com/game-jam-2026/dead-jump/internal/utils/audio"
)

//go:embed sfx/menu_music.wav
var menuMusicWAV []byte

//go:embed sfx/menu_select.wav
var menuSelectWAV []byte

//go:embed sfx/menu_confirm.wav
var menuConfirmWAV []byte

//go:embed sfx/cannonShot.mp3
var cannonShotMP3 []byte

//go:embed sfx/tsch.mp3
var projectileHitMP3 []byte

//go:embed sfx/spikeDeath1.mp3
var spikeDeath1MP3 []byte

//go:embed sfx/spikeDeath2.mp3
var spikeDeath2MP3 []byte

//go:embed sfx/gameOver.mp3
var gameOverMP3 []byte

func InitAudio() {
	audio.Init()

	_ = audio.RegisterMusicWAV(audio.SoundMenuMusic, menuMusicWAV)
	_ = audio.RegisterWAV(audio.SoundMenuSelect, menuSelectWAV)
	_ = audio.RegisterWAV(audio.SoundMenuConfirm, menuConfirmWAV)

	_ = audio.RegisterMP3(audio.SoundCannonShot, cannonShotMP3)
	_ = audio.RegisterMP3(audio.SoundProjectileHit, projectileHitMP3)

	_ = audio.RegisterMP3(audio.SoundDeath, spikeDeath1MP3)
	_ = audio.RegisterMP3(audio.SoundDeath, spikeDeath2MP3)
	_ = audio.RegisterMP3(audio.SoundGameOver, gameOverMP3)
}
