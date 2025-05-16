package raytracer

import (
	"math"
	"raytracer/math3"
)

type Sphere struct {
	Center   math3.Vec3
	Radius   float64
	Material Material
}

func (s Sphere) Hit(ray math3.Ray, rayT Interval, rec HitRecord) (bool, HitRecord) {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := math3.Dot(ray.Direction, oc)
	c := oc.LengthSquared() - s.Radius*s.Radius
	disc := h*h - a*c
	if disc < 0 {
		return false, rec
	}
	root := (h - math.Sqrt(disc)) / a
	if !rayT.Surrounds(root) {
		root = (h + math.Sqrt(disc)) / a
		if !rayT.Surrounds(root) {
			return false, rec
		}
	}
	rec.T = root
	rec.P = ray.At(rec.T)
	outwardNormal := rec.P.Sub(s.Center).Div(s.Radius)
	rec.SetFaceNormal(ray, outwardNormal)
	rec.Material = s.Material
	return true, rec
}
