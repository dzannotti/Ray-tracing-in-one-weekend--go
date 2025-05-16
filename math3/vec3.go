package math3

import "math"

const EPSILON = 0.0001

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (vec Vec3) Sub(other Vec3) Vec3 {
	return Vec3{
		X: vec.X - other.X,
		Y: vec.Y - other.Y,
		Z: vec.Z - other.Z,
	}
}

func (vec Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		X: vec.X + other.X,
		Y: vec.Y + other.Y,
		Z: vec.Z + other.Z,
	}
}

func (vec Vec3) Dot(other Vec3) float64 {
	return vec.X*other.X + vec.Y*other.Y + vec.Z*other.Z
}

func (vec Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: vec.Y*other.Z - vec.Z*other.Y,
		Y: vec.Z*other.X - vec.X*other.Z,
		Z: vec.X*other.Y - vec.Y*other.X,
	}
}

func (vec Vec3) Normalize() Vec3 {
	length := vec.Length()
	if length < EPSILON {
		panic("finding the length of an epsilon")
	}
	invLength := 1 / length
	return Vec3{
		X: vec.X * invLength,
		Y: vec.Y * invLength,
		Z: vec.Z * invLength,
	}
}

func (vec Vec3) K(n float64) Vec3 {
	return Vec3{
		X: vec.X * n,
		Y: vec.Y * n,
		Z: vec.Z * n,
	}
}

func (vec Vec3) Div(n float64) Vec3 {
	return vec.K(1 / n)
}

func (vec Vec3) LengthSquared() float64 {
	return vec.X*vec.X + vec.Y*vec.Y + vec.Z*vec.Z
}

func (vec Vec3) Length() float64 {
	return math.Sqrt(vec.LengthSquared())
}

func (vec Vec3) IsNearZero() bool {
	return math.Abs(vec.X) < EPSILON && math.Abs(vec.Y) < EPSILON && math.Abs(vec.Z) < EPSILON
}

func (vec Vec3) VectorMultiply(other Vec3) Vec3 {
	return Vec3{
		X: vec.X * other.X,
		Y: vec.Y * other.Y,
		Z: vec.Z * other.Z,
	}
}
