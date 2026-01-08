package noise

import (
	"image"
	"math"
	"math/rand"
)

func Gaussian(img *image.RGBA, deviation float64, perChannel bool) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			noiseR := rand.NormFloat64() * deviation

			noiseG := noiseR
			noiseB := noiseR

			if perChannel {
				noiseG = rand.NormFloat64() * deviation
				noiseB = rand.NormFloat64() * deviation
			}

			// noise RGB will have same values unless perChannel is enabled
			newR := clampFastFloat(float64(r) + noiseR)
			newG := clampFastFloat(float64(g) + noiseG)
			newB := clampFastFloat(float64(b) + noiseB)

			img.Pix[i+0] = newR
			img.Pix[i+1] = newG
			img.Pix[i+2] = newB
		}
	}
}

func clampFastFloat(val float64) uint8 {
	val = math.Round(val)
	// This achieves the same result as the if-statements
	return uint8(math.Max(0, math.Min(255, val)))
}
