package components

type Projectile struct {
	ImpulseMagnitude   float64
	DestroyOnHit       bool
	Lifetime           int
	MinSpeedForImpulse float64
	IsStationary       bool
}

type Cannon struct {
	FireRate            int
	FramesSinceLastShot int
	ProjectileSpeed     float64
	ProjectileMass      float64
	Direction           float64
	Active              bool
}

func DefaultCannon() Cannon {
	return Cannon{
		FireRate:            120,
		FramesSinceLastShot: 0,
		ProjectileSpeed:     4.0,
		ProjectileMass:      2.0,
		Direction:           3.14159,
		Active:              true,
	}
}

type ScreenSpace struct{}
