package linalg

import "math"

type Vector2 struct {
	X, Y float64
}

func Zero() Vector2 {
	return Vector2{0, 0}
}

func Up() Vector2 {
	return Vector2{0, -1}
}

func Down() Vector2 {
	return Vector2{0, 1}
}

func (v Vector2) Add(b Vector2) Vector2 {
	return Vector2{v.X + b.X, v.Y + b.Y}
}

func (v Vector2) Sub(b Vector2) Vector2 {
	return Vector2{v.X - b.X, v.Y - b.Y}
}

func (v Vector2) Scale(s float64) Vector2 {
	return Vector2{v.X * s, v.Y * s}
}

func (v Vector2) Mul(b Vector2) Vector2 {
	return Vector2{v.X * b.X, v.Y * b.Y}
}

func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) Normalized() Vector2 {
	length := v.Length()
	if length == 0 {
		return Zero()
	}
	return v.Scale(1 / length)
}

func (v Vector2) Dot(b Vector2) float64 {
	return v.X*b.X + v.Y*b.Y
}

func (v Vector2) Lerp(b Vector2, t float64) Vector2 {
	return Vector2{
		X: v.X + (b.X-v.X)*t,
		Y: v.Y + (b.Y-v.Y)*t,
	}
}

func (v Vector2) Clamp(min, max float64) Vector2 {
	return Vector2{
		X: math.Max(min, math.Min(max, v.X)),
		Y: math.Max(min, math.Min(max, v.Y)),
	}
}

func (v Vector2) ClampLength(maxLength float64) Vector2 {
	length := v.Length()
	if length > maxLength && length > 0 {
		return v.Scale(maxLength / length)
	}
	return v
}

func (v Vector2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vector2) Perpendicular() Vector2 {
	return Vector2{-v.Y, v.X}
}

func (v Vector2) Reflect(normal Vector2) Vector2 {
	dot := v.Dot(normal)
	return v.Sub(normal.Scale(2 * dot))
}

func (v Vector2) Project(onto Vector2) Vector2 {
	dot := v.Dot(onto)
	lengthSq := onto.LengthSquared()
	if lengthSq == 0 {
		return Zero()
	}
	return onto.Scale(dot / lengthSq)
}

func (v Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func FromAngle(radians float64) Vector2 {
	return Vector2{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}
