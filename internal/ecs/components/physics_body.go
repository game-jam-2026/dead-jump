package components

import "github.com/game-jam-2026/dead-jump/pkg/linalg"

type PhysicsBody struct {
	Mass         float64
	Friction     float64
	Bounciness   float64
	AirDrag      float64
	GravityScale float64
	IsKinematic  bool
	IsGrounded   bool
	GroundNormal linalg.Vector2
	MaxSpeed     float64
	Acceleration linalg.Vector2
}

func DefaultPhysicsBody() PhysicsBody {
	return PhysicsBody{
		Mass:         1.0,
		Friction:     0.3,
		Bounciness:   0.0,
		AirDrag:      0.02,
		GravityScale: 1.0,
		IsKinematic:  false,
		IsGrounded:   false,
		GroundNormal: linalg.Up(),
		MaxSpeed:     10.0,
		Acceleration: linalg.Zero(),
	}
}

func StaticBody() PhysicsBody {
	return PhysicsBody{
		Mass:         0,
		Friction:     0.5,
		Bounciness:   0.0,
		AirDrag:      0.0,
		GravityScale: 0.0,
		IsKinematic:  true,
		IsGrounded:   true,
		GroundNormal: linalg.Up(),
		MaxSpeed:     0,
		Acceleration: linalg.Zero(),
	}
}

func ProjectileBody(mass float64) PhysicsBody {
	return PhysicsBody{
		Mass:         mass,
		Friction:     0.1,
		Bounciness:   0.3,
		AirDrag:      0.01,
		GravityScale: 1.0,
		IsKinematic:  false,
		IsGrounded:   false,
		GroundNormal: linalg.Zero(),
		MaxSpeed:     20.0,
		Acceleration: linalg.Zero(),
	}
}

func (pb *PhysicsBody) AddForce(force linalg.Vector2) {
	if pb.Mass > 0 {
		pb.Acceleration = pb.Acceleration.Add(force.Scale(1 / pb.Mass))
	}
}

func (pb *PhysicsBody) AddImpulse(impulse linalg.Vector2) linalg.Vector2 {
	if pb.Mass > 0 {
		return impulse.Scale(1 / pb.Mass)
	}
	return linalg.Zero()
}

func (pb *PhysicsBody) IsStatic() bool {
	return pb.Mass <= 0 || pb.IsKinematic
}
