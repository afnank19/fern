package point

import (
	"image"
	"math"
)

func LinearContrast(img *image.RGBA, factor float64) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			newR := factor*(float64(r)-128.0) + 128.0
			newG := factor*(float64(g)-128.0) + 128.0
			newB := factor*(float64(b)-128.0) + 128.0

			img.Pix[i+0] = clampFastFloat(newR)
			img.Pix[i+1] = clampFastFloat(newG)
			img.Pix[i+2] = clampFastFloat(newB)
		}
	}
}

// A lot of code repitition here, refactor into smaller chunks
// This code can also be optimized a lot, look into LUTs, lowBound, highBound requirements etc
func SigmoidalContrast(img *image.RGBA, factor float64) {
	if factor == 0 {
		return
	}

	bounds := img.Bounds()

	// optimizations can be done here using Look Up tables
	lowBound := sigmoid(factor, 0)
	highBound := sigmoid(factor, 1)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			normR := float64(r) / 255.0
			normG := float64(g) / 255.0
			normB := float64(b) / 255.0

			rSig := sigmoid(factor, normR)
			gSig := sigmoid(factor, normG)
			bSig := sigmoid(factor, normB)

			newR := (rSig - lowBound) / (highBound - lowBound)
			newG := (gSig - lowBound) / (highBound - lowBound)
			newB := (bSig - lowBound) / (highBound - lowBound)

			img.Pix[i+0] = clampFastFloat(newR * 255.0)
			img.Pix[i+1] = clampFastFloat(newG * 255.0)
			img.Pix[i+2] = clampFastFloat(newB * 255.0)
		}
	}
}

func sigmoid(a, b float64) float64 {
	return 1.0 / (1.0 + math.Exp(-a*(b-0.5)))
}

func FastSigmoidalContrast(img *image.RGBA, factor float64) {
	if factor == 0 {
		return
	}

	bounds := img.Bounds()

	lowBound := sigmoid(factor, 0)
	highBound := sigmoid(factor, 1)
	lut := generateSigmoidalLUT(factor, lowBound, highBound)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			newR := lut[r]
			newG := lut[g]
			newB := lut[b]

			img.Pix[i+0] = clampFastFloat(newR * 255.0)
			img.Pix[i+1] = clampFastFloat(newG * 255.0)
			img.Pix[i+2] = clampFastFloat(newB * 255.0)
		}
	}
}

func generateSigmoidalLUT(factor, lB, hB float64) []float64 {
	lut := make([]float64, 256)
	for i := range lut {
		normI := float64(i) / 255.0

		y := sigmoid(factor, normI)
		newY := (y - lB) / (hB - lB)
		lut[i] = newY
	}

	return lut
}

// might move later to utils
func clampFloat(value float64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

func clampFastFloat(val float64) uint8 {
	// This achieves the same result as the if-statements
	return uint8(math.Max(0, math.Min(255, val)))
}
