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
