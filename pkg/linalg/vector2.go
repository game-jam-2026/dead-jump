package linalg

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Add(b Vector2) Vector2 {
	return Vector2{v.X + b.X, v.Y + b.Y}
}

func (v Vector2) Scale(s float64) Vector2 {
	return Vector2{v.X * s, v.Y * s}
}
