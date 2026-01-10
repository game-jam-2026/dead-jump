package audio

import (
	"bytes"
	"io"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Manager struct {
	context      *audio.Context
	sounds       map[SoundID][][]byte
	musicPlayers map[SoundID]*audio.Player
	muted        bool
	masterVolume float64
	musicVolume  float64
	sfxVolume    float64
	initialized  bool
}

var defaultManager *Manager

func Init() {
	if defaultManager != nil && defaultManager.initialized {
		return
	}
	defaultManager = &Manager{
		context:      audio.NewContext(44100),
		sounds:       make(map[SoundID][][]byte),
		musicPlayers: make(map[SoundID]*audio.Player),
		masterVolume: 1.0,
		musicVolume:  0.7,
		sfxVolume:    1.0,
		initialized:  true,
	}
}

func GetContext() *audio.Context {
	if defaultManager == nil {
		return nil
	}
	return defaultManager.context
}

func RegisterWAV(id SoundID, data []byte) error {
	if defaultManager == nil {
		Init()
	}
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	decoded, err := io.ReadAll(stream)
	if err != nil {
		return err
	}
	defaultManager.sounds[id] = append(defaultManager.sounds[id], decoded)
	return nil
}

func RegisterMP3(id SoundID, data []byte) error {
	if defaultManager == nil {
		Init()
	}
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	decoded, err := io.ReadAll(stream)
	if err != nil {
		return err
	}
	defaultManager.sounds[id] = append(defaultManager.sounds[id], decoded)
	return nil
}

func RegisterMusicWAV(id SoundID, data []byte) error {
	if defaultManager == nil {
		Init()
	}
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := defaultManager.context.NewPlayer(loop)
	if err != nil {
		return err
	}
	defaultManager.musicPlayers[id] = player
	return nil
}

func Play(id SoundID) {
	if defaultManager == nil || defaultManager.muted {
		return
	}
	variants := defaultManager.sounds[id]
	if len(variants) == 0 {
		return
	}
	sound := variants[rand.Intn(len(variants))]
	player := defaultManager.context.NewPlayerFromBytes(sound)
	player.SetVolume(GetEffectiveSFXVolume())
	player.Play()
}

func PlayMusic(id SoundID) {
	if defaultManager == nil {
		return
	}
	player := defaultManager.musicPlayers[id]
	if player == nil {
		return
	}
	if !defaultManager.muted && !player.IsPlaying() {
		player.SetVolume(GetEffectiveMusicVolume())
		player.Play()
	}
}

func StopMusic(id SoundID) {
	if defaultManager == nil {
		return
	}
	player := defaultManager.musicPlayers[id]
	if player != nil && player.IsPlaying() {
		player.Pause()
	}
}

func IsMusicPlaying(id SoundID) bool {
	if defaultManager == nil {
		return false
	}
	player := defaultManager.musicPlayers[id]
	return player != nil && player.IsPlaying()
}

func UpdateMusicVolume() {
	if defaultManager == nil {
		return
	}
	vol := GetEffectiveMusicVolume()
	for _, player := range defaultManager.musicPlayers {
		if player != nil {
			player.SetVolume(vol)
		}
	}
}

func IsMuted() bool {
	return defaultManager != nil && defaultManager.muted
}

func SetMuted(muted bool) {
	if defaultManager == nil {
		return
	}
	defaultManager.muted = muted
	if muted {
		for _, player := range defaultManager.musicPlayers {
			if player != nil && player.IsPlaying() {
				player.Pause()
			}
		}
	}
}

func GetMasterVolume() float64 {
	if defaultManager == nil {
		return 1.0
	}
	return defaultManager.masterVolume
}

func SetMasterVolume(v float64) {
	if defaultManager == nil {
		return
	}
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	defaultManager.masterVolume = v
}

func GetMusicVolume() float64 {
	if defaultManager == nil {
		return 0.7
	}
	return defaultManager.musicVolume
}

func SetMusicVolume(v float64) {
	if defaultManager == nil {
		return
	}
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	defaultManager.musicVolume = v
}

func GetSFXVolume() float64 {
	if defaultManager == nil {
		return 1.0
	}
	return defaultManager.sfxVolume
}

func SetSFXVolume(v float64) {
	if defaultManager == nil {
		return
	}
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	defaultManager.sfxVolume = v
}

func GetEffectiveSFXVolume() float64 {
	if defaultManager == nil {
		return 1.0
	}
	return defaultManager.sfxVolume * defaultManager.masterVolume
}

func GetEffectiveMusicVolume() float64 {
	if defaultManager == nil {
		return 0.7
	}
	return defaultManager.musicVolume * defaultManager.masterVolume
}
