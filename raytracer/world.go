package raytracer

import (
	"raytracer/math3"
)

type World struct {
	Objects []Hittable
}

func (w *World) Clear() {
	w.Objects = make([]Hittable, 0)
}

func (w *World) Add(obj Hittable) {
	w.Objects = append(w.Objects, obj)
}

func (w *World) Hit(ray math3.Ray, rayT Interval) (HitRecord, bool) {
	hitAnything := false
	closestSoFar := rayT.Max
	rec := HitRecord{}
	for _, obj := range w.Objects {
		localRec, hasHit := obj.Hit(ray, Interval{Min: rayT.Min, Max: closestSoFar})
		if hasHit {
			hitAnything = true
			closestSoFar = localRec.T
			rec = localRec
		}
	}
	return rec, hitAnything
}
