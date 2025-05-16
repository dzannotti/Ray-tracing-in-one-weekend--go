package raytracer

import (
	"image/color"
	"math"
	"raytracer/math3"
)

func colorToInt(val float64) uint8 {
	return uint8(math.Floor(val))
}

func linearToGamma(val float64) float64 {
	if val > 0 {
		return math.Sqrt(val)
	}
	return 0
}

func convertPixel(pixel math3.Vec3) color.Color {
	return color.RGBA{
		R: colorToInt(linearToGamma(pixel[0]) * 255),
		G: colorToInt(linearToGamma(pixel[1]) * 255),
		B: colorToInt(linearToGamma(pixel[2]) * 255),
		A: 255,
	}
}
