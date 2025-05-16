package math3

import (
	"math"
	"math/rand/v2"
)

func Deg2Rad(deg float64) float64 {
	return (deg / 180.0) * math.Pi
}

func Dot(v Vec3, u Vec3) float64 {
	return v.Dot(u)
}

func Reflect(v Vec3, n Vec3) Vec3 {
	return v.Sub(n.Scale(v.Dot(n) * 2))
}

func Random() Vec3 {
	return Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func RandomBetween(low float64, high float64) Vec3 {
	return Vec3{
		low + (high-low)*rand.Float64(),
		low + (high-low)*rand.Float64(),
		low + (high-low)*rand.Float64(),
	}
}

func RandomUnitVector() Vec3 {
	for {
		p := RandomBetween(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1.0 {
			return p.Div(math.Sqrt(lensq))
		}
	}
}

func RandomOnHemisphere(normal Vec3) Vec3 {
	onSphere := RandomUnitVector()
	if Dot(onSphere, normal) > 0 {
		return onSphere
	}
	return onSphere.Scale(-1)
}

func RandomInUnitDisk() Vec3 {
	for {
		p := RandomBetween(-1, 1)
		p[2] = 0
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func Refract(uv Vec3, n Vec3, etaiOverEtat float64) Vec3 {
	cosT := math.Min(Dot(uv.Scale(-1), n), 1.0)
	rOutPerp := uv.Add(n.Scale(cosT)).Scale(etaiOverEtat)
	rOutParallel := n.Scale(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func Cross(u Vec3, v Vec3) Vec3 {
	return u.Cross(v)
}
