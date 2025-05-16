package raytracer

import (
	"image"
	"image/color"
	"math"
)

func BilateralFilter(input image.Image, spatialSigma, rangeSigma float64) *image.RGBA {
	bounds := input.Bounds()
	output := image.NewRGBA(bounds)

	// Calculate kernel radius based on spatialSigma
	// Use 3 * sigma to capture most of the Gaussian curve
	kernelRadius := int(math.Ceil(3.0 * spatialSigma))

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get center pixel color
			centerR, centerG, centerB, centerA := getRGBAFloats(input.At(x, y))

			// Accumulators for filtered color and weights
			var sumR, sumG, sumB, sumA, totalWeight float64

			// Apply filter kernel
			for ky := -kernelRadius; ky <= kernelRadius; ky++ {
				for kx := -kernelRadius; kx <= kernelRadius; kx++ {
					// Neighbor coordinates
					nx, ny := x+kx, y+ky

					// Skip if outside image bounds
					if nx < bounds.Min.X || nx >= bounds.Max.X || ny < bounds.Min.Y || ny >= bounds.Max.Y {
						continue
					}

					// Get neighbor pixel color
					neighborR, neighborG, neighborB, neighborA := getRGBAFloats(input.At(nx, ny))

					// Calculate spatial weight (based on distance)
					spatialDist := float64(kx*kx + ky*ky)
					spatialWeight := math.Exp(-spatialDist / (2.0 * spatialSigma * spatialSigma))

					// Calculate range weight (based on color difference)
					colorDist := colorDistance(
						centerR, centerG, centerB, centerA,
						neighborR, neighborG, neighborB, neighborA,
					)
					rangeWeight := math.Exp(-colorDist / (2.0 * rangeSigma * rangeSigma))

					// Combined weight
					weight := spatialWeight * rangeWeight

					// Accumulate weighted color
					sumR += neighborR * weight
					sumG += neighborG * weight
					sumB += neighborB * weight
					sumA += neighborA * weight
					totalWeight += weight
				}
			}

			// Normalize by total weight
			if totalWeight > 0 {
				sumR /= totalWeight
				sumG /= totalWeight
				sumB /= totalWeight
				sumA /= totalWeight
			}

			// Set output pixel
			output.Set(x, y, color.RGBA{
				R: uint8(math.Min(math.Max(0, sumR*255), 255)),
				G: uint8(math.Min(math.Max(0, sumG*255), 255)),
				B: uint8(math.Min(math.Max(0, sumB*255), 255)),
				A: uint8(math.Min(math.Max(0, sumA*255), 255)),
			})
		}
	}

	return output
}

func getRGBAFloats(c color.Color) (r, g, b, a float64) {
	r32, g32, b32, a32 := c.RGBA()
	return float64(r32) / 65535.0, float64(g32) / 65535.0, float64(b32) / 65535.0, float64(a32) / 65535.0
}

func colorDistance(r1, g1, b1, a1, r2, g2, b2, a2 float64) float64 {
	dr := r1 - r2
	dg := g1 - g2
	db := b1 - b2
	da := a1 - a2
	return dr*dr + dg*dg + db*db + da*da
}
