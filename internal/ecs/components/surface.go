package components

import (
	"math"

	"github.com/game-jam-2026/dead-jump/pkg/linalg"
)

type SurfaceType int

const (
	SurfaceNormal SurfaceType = iota
	SurfaceIce
	SurfaceRough
	SurfaceSticky
	SurfaceConveyor
	SurfaceBouncy
)

type Surface struct {
	Type              SurfaceType
	Friction          float64
	Bounciness        float64
	SlopeAngle        float64
	Normal            linalg.Vector2
	ConveyorSpeed     float64
	ConveyorDirection linalg.Vector2
	IsPlatform        bool
	Tags              []string
}

func NewSurface(surfaceType SurfaceType) Surface {
	s := Surface{
		Type:              surfaceType,
		Normal:            linalg.Up(),
		ConveyorDirection: linalg.Vector2{X: 1, Y: 0},
	}

	switch surfaceType {
	case SurfaceNormal:
		s.Friction = 0.5
		s.Bounciness = 0.0
	case SurfaceIce:
		s.Friction = 0.05
		s.Bounciness = 0.0
	case SurfaceRough:
		s.Friction = 0.8
		s.Bounciness = 0.0
	case SurfaceSticky:
		s.Friction = 0.95
		s.Bounciness = 0.0
	case SurfaceConveyor:
		s.Friction = 0.4
		s.Bounciness = 0.0
		s.ConveyorSpeed = 2.0
	case SurfaceBouncy:
		s.Friction = 0.3
		s.Bounciness = 0.8
	}

	return s
}

func NewSlopedSurface(surfaceType SurfaceType, slopeAngle float64) Surface {
	s := NewSurface(surfaceType)
	s.SlopeAngle = slopeAngle
	s.Normal = linalg.Vector2{
		X: -math.Sin(slopeAngle),
		Y: -math.Cos(slopeAngle),
	}
	return s
}

func NewConveyor(speed float64, direction linalg.Vector2) Surface {
	s := NewSurface(SurfaceConveyor)
	s.ConveyorSpeed = speed
	s.ConveyorDirection = direction.Normalized()
	return s
}

func (s *Surface) GetEffectiveFriction(bodyFriction float64) float64 {
	return math.Sqrt(s.Friction * bodyFriction)
}

func (s *Surface) GetSlopeDirection() linalg.Vector2 {
	return linalg.Vector2{
		X: math.Cos(s.SlopeAngle),
		Y: math.Sin(s.SlopeAngle),
	}
}

func (s *Surface) GetSlideAcceleration(gravity float64) float64 {
	return gravity * math.Sin(s.SlopeAngle)
}
