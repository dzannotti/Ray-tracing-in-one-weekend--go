package raytracer

import (
	"math"
	"math/rand/v2"
	"raytracer/math3"
)

type Material interface {
	Scatter(ray math3.Ray, rec HitRecord) (math3.Vec3, math3.Ray, bool)
}

type Lambertian struct {
	Albedo math3.Vec3
}

func (l Lambertian) Scatter(ray math3.Ray, rec HitRecord) (math3.Vec3, math3.Ray, bool) {
	scatterDir := rec.Normal.Add(math3.RandomUnitVector())
	if scatterDir.IsNearZero() {
		scatterDir = rec.Normal
	}
	return l.Albedo, math3.Ray{Origin: rec.P, Direction: scatterDir}, true
}

type Metal struct {
	Fuzz   float64
	Albedo math3.Vec3
}

func (m Metal) Scatter(ray math3.Ray, rec HitRecord) (math3.Vec3, math3.Ray, bool) {
	reflected := math3.Reflect(ray.Direction, rec.Normal)
	reflected = reflected.Normalize().Add(math3.RandomUnitVector().Scale(m.Fuzz))
	scattered := math3.Ray{Origin: rec.P, Direction: reflected}
	canScatter := math3.Dot(scattered.Direction, rec.Normal) > 0
	return m.Albedo, scattered, canScatter
}

type Dialectric struct {
	RefractionIndex float64
}

func (d Dialectric) Scatter(ray math3.Ray, rec HitRecord) (math3.Vec3, math3.Ray, bool) {
	ri := d.RefractionIndex
	if rec.FrontFace {
		ri = 1 / d.RefractionIndex
	}
	unitDir := ray.Direction.Normalize()
	cosT := math.Min(math3.Dot(unitDir.Scale(-1), rec.Normal), 1)
	sinT := math.Sqrt(math.Max(0.0, 1.0-cosT*cosT))
	cannotRefract := ri*sinT > 1
	var direction math3.Vec3
	if cannotRefract || d.reflectance(cosT, ri) > rand.Float64() {
		direction = math3.Reflect(unitDir, rec.Normal)
	} else {
		direction = math3.Refract(unitDir, rec.Normal, ri)
	}
	scattered := math3.Ray{Origin: rec.P, Direction: direction}
	return math3.Vec3{1, 1, 1}, scattered, true
}

func (d Dialectric) reflectance(cosine float64, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
