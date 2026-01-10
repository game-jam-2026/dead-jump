package physics

import "github.com/game-jam-2026/dead-jump/pkg/linalg"

type Config struct {
	Gravity             linalg.Vector2
	TerminalVelocity    float64
	DefaultFriction     float64
	MinVelocity         float64
	FixedDeltaTime      float64
	CollisionIterations int
	PushForce           float64
	SlopeThreshold      float64
	GroundCheckDistance float64
	MaxStepDistance     float64
}

func DefaultConfig() *Config {
	return &Config{
		Gravity:             linalg.Vector2{X: 0, Y: 0.5},
		TerminalVelocity:    15.0,
		DefaultFriction:     0.5,
		MinVelocity:         0.01,
		FixedDeltaTime:      1.0 / 60.0,
		CollisionIterations: 4,
		PushForce:           5.0,
		SlopeThreshold:      0.785,
		GroundCheckDistance: 2.0,
		MaxStepDistance:     8.0,
	}
}
