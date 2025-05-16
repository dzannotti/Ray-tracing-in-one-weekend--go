package raytracer

import "raytracer/math3"

type World struct {
	Objects []Hittable
}

func (w *World) Clear() {
	w.Objects = make([]Hittable, 0)
}

func (w *World) Add(obj Hittable) {
	w.Objects = append(w.Objects, obj)
}

func (w *World) Hit(ray math3.Ray, rayT Interval, rec HitRecord) (bool, HitRecord) {
	tempRec := HitRecord{}
	hitAnything := false
	closestSoFar := rayT.Max
	for _, obj := range w.Objects {
		hasHit, resultRec := obj.Hit(ray, Interval{Min: rayT.Min, Max: closestSoFar}, tempRec)
		if hasHit {
			hitAnything = true
			closestSoFar = resultRec.T
			rec = resultRec
		}
	}
	return hitAnything, rec
}
