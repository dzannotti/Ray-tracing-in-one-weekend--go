package math3

import "math"

const EPSILON = 1e-8

type Vec3 [3]float64

func (vec Vec3) X() float64 { return vec[0] }
func (vec Vec3) Y() float64 { return vec[1] }
func (vec Vec3) Z() float64 { return vec[2] }

func (vec Vec3) Sub(other Vec3) Vec3 {
	return Vec3{
		vec[0] - other[0],
		vec[1] - other[1],
		vec[2] - other[2],
	}
}

func (vec Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		vec[0] + other[0],
		vec[1] + other[1],
		vec[2] + other[2],
	}
}

func (vec Vec3) Scale(k float64) Vec3 {
	return Vec3{
		vec[0] * k,
		vec[1] * k,
		vec[2] * k,
	}
}

func (vec Vec3) Div(k float64) Vec3 {
	return vec.Scale(1 / k)
}

func (vec Vec3) Dot(other Vec3) float64 {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2]
}

func (vec Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		vec[1]*other[2] - vec[2]*other[1],
		vec[2]*other[0] - vec[0]*other[2],
		vec[0]*other[1] - vec[1]*other[0],
	}
}

func (vec Vec3) Normalize() Vec3 {
	sl := vec.LengthSquared()
	if sl < EPSILON {
		return vec
	}
	return vec.Scale(1 / math.Sqrt(sl))
}

func (vec Vec3) Length() float64 {
	return math.Sqrt(vec.LengthSquared())
}

func (vec Vec3) LengthSquared() float64 {
	return vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]
}

func (vec Vec3) IsNearZero() bool {
	return math.Abs(vec[0]) < EPSILON && math.Abs(vec[1]) < EPSILON && math.Abs(float64(vec[2])) < EPSILON
}

func (vec Vec3) Multiply(v Vec3) Vec3 {
	return Vec3{vec[0] * v[0], vec[1] * v[1], vec[2] * v[2]}
}

func (v Vec3) MaxComponent() float64 {
	return math.Max(math.Max(v.X(), v.Y()), v.Z())
}
