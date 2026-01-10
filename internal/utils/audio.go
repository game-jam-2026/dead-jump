package utils

import (
	"bytes"
	"math/rand"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
)

func PlaySound(w *ecs.World, sound []byte) error {
	audioEntities := w.GetEntities(
		reflect.TypeOf((*components.AudioContext)(nil)).Elem(),
	)

	if len(audioEntities) == 0 {
		return nil
	}

	audioCtx, err := ecs.GetComponent[components.AudioContext](w, audioEntities[0])
	if err != nil {
		return err
	}

	stream, err := mp3.DecodeWithSampleRate(audioCtx.Context.SampleRate(), bytes.NewReader(sound))
	if err != nil {
		return err
	}

	player, err := audio.NewPlayer(audioCtx.Context, stream)
	if err != nil {
		return err
	}

	player.Play()
	return nil
}

func PlayRandomSound(w *ecs.World, sounds [][]byte) error {
	if len(sounds) == 0 {
		return nil
	}
	idx := rand.Intn(len(sounds))
	return PlaySound(w, sounds[idx])
}
