package assets

import (
	"bytes"
	_ "embed"
	"io"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed sfx/cannonShot.mp3
var CannonShotMP3 []byte

//go:embed sfx/tsch.mp3
var TschMP3 []byte

//go:embed sfx/spikeDeath1.mp3
var SpikeDeath1MP3 []byte

//go:embed sfx/spikeDeath2.mp3
var SpikeDeath2MP3 []byte

var (
	AudioContext     *audio.Context
	CannonShotBytes  []byte
	TschBytes        []byte
	SpikeDeath1Bytes []byte
	SpikeDeath2Bytes []byte
	SpikeDeathSounds [][]byte
)

func InitAudio(ctx *audio.Context) {
	AudioContext = ctx

	shotDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(CannonShotMP3))
	if err != nil {
		panic(err)
	}

	CannonShotBytes, err = io.ReadAll(shotDecoded)
	if err != nil {
		panic(err)
	}

	tschDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(TschMP3))
	if err != nil {
		panic(err)
	}

	TschBytes, err = io.ReadAll(tschDecoded)
	if err != nil {
		panic(err)
	}

	death1Decoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(SpikeDeath1MP3))
	if err != nil {
		panic(err)
	}

	SpikeDeath1Bytes, err = io.ReadAll(death1Decoded)
	if err != nil {
		panic(err)
	}

	death2Decoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(SpikeDeath2MP3))
	if err != nil {
		panic(err)
	}

	SpikeDeath2Bytes, err = io.ReadAll(death2Decoded)
	if err != nil {
		panic(err)
	}

	SpikeDeathSounds = [][]byte{SpikeDeath1Bytes, SpikeDeath2Bytes}
}

func PlayCannonShot() {
	player := AudioContext.NewPlayerFromBytes(CannonShotBytes)
	player.Play()
}

func PlayTsch() {
	player := AudioContext.NewPlayerFromBytes(TschBytes)
	player.Play()
}

func PlayRandomDeathSound() {
	if len(SpikeDeathSounds) == 0 {
		return
	}
	sound := SpikeDeathSounds[rand.Intn(len(SpikeDeathSounds))]
	player := AudioContext.NewPlayerFromBytes(sound)
	player.Play()
}
