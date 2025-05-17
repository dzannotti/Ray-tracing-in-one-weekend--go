package raytracer

import "raytracer/math3"

type Hittable interface {
	Origin() math3.Vec3
	Hit(ray math3.Ray, rayT Interval) (HitRecord, bool)
	Prepare()
}

type HitRecord struct {
	FrontFace bool
	P         math3.Vec3
	Normal    math3.Vec3
	T         float64
	Material  Material
}

func (hr *HitRecord) SetFaceNormal(r math3.Ray, outwardNormal math3.Vec3) {
	hr.FrontFace = math3.Dot(r.Direction, outwardNormal) < 0
	hr.Normal = outwardNormal
	if !hr.FrontFace {
		hr.Normal = hr.Normal.Scale(-1)
	}
}
