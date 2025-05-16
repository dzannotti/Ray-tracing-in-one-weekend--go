package raytracer

import (
	"raytracer/math3"
)

type Material interface {
	Scatter(rayIn **math3.Ray, rec HitRecord, attenuation **math3.Vec3, scattered **math3.Ray) bool
}

type Lambertian struct {
	Albedo math3.Vec3
}

func (l Lambertian) Scatter(_ **math3.Ray, rec HitRecord, attenuation **math3.Vec3, scattered **math3.Ray) bool {
	scatterDir := rec.Normal.Add(math3.RandomUnitVector())
	if scatterDir.IsNearZero() {
		scatterDir = rec.Normal
	}
	*scattered = &math3.Ray{Origin: rec.P, Direction: scatterDir}
	*attenuation = &l.Albedo
	return true
}

// TODO: port the other materials
