package raytracer

import (
	"raytracer/math3"
	"slices"
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

func (w *World) Prepare() {
	for _, obj := range w.Objects {
		obj.Prepare()
	}
	// closest first should lead to less iterations
	slices.SortFunc(w.Objects, func(a Hittable, b Hittable) int {
		return int(a.Origin().Z() - b.Origin().Z())
	})

}

func (w *World) Hit(ray math3.Ray, rayT Interval) (HitRecord, bool) {
	hitAnything := false
	closestSoFar := rayT.Max
	rec := HitRecord{}
	for _, obj := range w.Objects {
		if localRec, hasHit := obj.Hit(ray, Interval{Min: rayT.Min, Max: closestSoFar}); hasHit {
			hitAnything = true
			closestSoFar = localRec.T
			rec = localRec
		}
	}
	return rec, hitAnything
}
