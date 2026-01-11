package audio

type SoundID int

const (
	SoundMenuMusic SoundID = iota
	SoundMenuSelect
	SoundMenuConfirm

	SoundCannonShot
	SoundProjectileHit
	SoundDeath
	SoundGameOver

	SoundLevelMusic
	SoundStep
)
