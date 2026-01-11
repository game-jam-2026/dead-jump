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

	BurstCount           int
	BurstDelay           int
	CurrentBurst         int
	FramesSinceLastBurst int
}

func DefaultCannon() Cannon {
	return Cannon{
		FireRate:             120,
		FramesSinceLastShot:  0,
		ProjectileSpeed:      4.0,
		ProjectileMass:       2.0,
		Direction:            3.14159,
		Active:               true,
		BurstCount:           1,
		BurstDelay:           5,
		CurrentBurst:         0,
		FramesSinceLastBurst: 0,
	}
}

type ScreenSpace struct{}
