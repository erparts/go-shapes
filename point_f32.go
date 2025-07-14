package shapes

import "math"

type PointF32 struct {
	X, Y float32
}

func (p PointF32) Sub(o PointF32) PointF32 {
	return PointF32{p.X - o.X, p.Y - o.Y}
}

func (p PointF32) Add(o PointF32) PointF32 {
	return PointF32{p.X + o.X, p.Y + o.Y}
}

func (p PointF32) Scale(s float32) PointF32 {
	return PointF32{p.X * s, p.Y * s}
}

func (p PointF32) Dot(o PointF32) float32 {
	return p.X*o.X + p.Y*o.Y
}

func (p PointF32) Length() float32 {
	return float32(math.Hypot(float64(p.X), float64(p.Y)))
}

func (p PointF32) Normalize() PointF32 {
	l := p.Length()
	if l == 0 {
		return PointF32{0, 0}
	}
	return PointF32{p.X / l, p.Y / l}
}
