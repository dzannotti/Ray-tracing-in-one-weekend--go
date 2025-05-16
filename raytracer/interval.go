package raytracer

import "math"

type Interval struct {
	Min float64
	Max float64
}

var (
	EmptyInterval    = Interval{Min: math.MaxFloat64, Max: -math.MaxFloat64}
	UniverseInterval = Interval{Min: -math.MaxFloat64, Max: math.MaxFloat64}
)

func (iv Interval) Size() float64 {
	return iv.Max - iv.Min
}

func (iv Interval) Contains(v float64) bool {
	return iv.Min <= v && v <= iv.Max
}

func (iv Interval) Surrounds(v float64) bool {
	return iv.Min < v && v < iv.Max
}

func (iv Interval) Clamp(v float64) float64 {
	if v < iv.Min {
		return iv.Min
	}
	if v > iv.Max {
		return iv.Max
	}
	return v
}
