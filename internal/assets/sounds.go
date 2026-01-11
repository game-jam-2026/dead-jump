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

//go:embed music/flex.mp3
var levelMusicMP3 []byte

//go:embed sfx/steps/step1.mp3
var step1MP3 []byte

//go:embed sfx/steps/step2.mp3
var step2MP3 []byte

//go:embed sfx/steps/step3.mp3
var step3MP3 []byte

//go:embed sfx/steps/step4.mp3
var step4MP3 []byte

//go:embed sfx/steps/step5.mp3
var step5MP3 []byte

//go:embed sfx/steps/step6.mp3
var step6MP3 []byte

//go:embed sfx/steps/step7.mp3
var step7MP3 []byte

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

	_ = audio.RegisterMusicMP3(audio.SoundLevelMusic, levelMusicMP3)

	_ = audio.RegisterMP3(audio.SoundStep, step1MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step2MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step3MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step4MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step5MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step6MP3)
	_ = audio.RegisterMP3(audio.SoundStep, step7MP3)
}
