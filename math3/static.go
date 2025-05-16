package math3

import (
	"math"
	"math/rand/v2"
)

const Ded2Rad = (1 / 180.0) * math.Pi

func Dot(v Vec3, u Vec3) float64 {
	return v.Dot(u)
}

func Reflect(v Vec3, n Vec3) Vec3 {
	return v.Sub(n.K(Dot(v, n) * 2))
}

func Random() Vec3 {
	return Vec3{
		X: rand.Float64(),
		Y: rand.Float64(),
		Z: rand.Float64(),
	}
}

func RandomBetween(low float64, high float64) Vec3 {
	return Vec3{
		X: low + rand.Float64()*(high-low),
		Y: low + rand.Float64()*(high-low),
		Z: low + rand.Float64()*(high-low),
	}
}

func RandomUnitVector() Vec3 {
	for {
		p := RandomBetween(-1, 1)
		l := p.LengthSquared()
		if EPSILON < l && l <= 1 {
			return p.Normalize()
		}
	}
}

func RandomOnHemisphere(normal Vec3) Vec3 {
	onSphere := RandomUnitVector()
	if Dot(onSphere, normal) > 0 {
		return onSphere
	}
	return onSphere.K(-1)
}

func RandomInUnitDisk() Vec3 {
	for {
		p := RandomBetween(-1, 1)
		p.Z = 0
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func Refract(uv Vec3, n Vec3, etaiOverEtat float64) Vec3 {
	cosT := math.Min(Dot(uv.K(-1), n), 1)
	rOutPerp := uv.Add(n.K(cosT).K(etaiOverEtat))
	rOutParallel := n.K(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func Cross(u Vec3, v Vec3) Vec3 {
	return u.Cross(v)
}
