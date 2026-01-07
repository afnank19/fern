package point

import (
	"fmt"
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
func SigmoidalContrast(img *image.RGBA, factor float64) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			normR := float64(r) / 255.0
			normG := float64(g) / 255.0
			normB := float64(b) / 255.0

			rSig := 1.0 / (1.0 + math.Exp(-factor*(normR-0.5)))
			gSig := 1.0 / (1.0 + math.Exp(-factor*(normG-0.5)))
			bSig := 1.0 / (1.0 + math.Exp(-factor*(normB-0.5)))

			// optimizations can be done here using Look Up tables
			lowBound := 1.0 / (1.0 + math.Exp(-factor*(0-0.5)))
			highBound := 1.0 / (1.0 + math.Exp(-factor*(1-0.5)))

			newR := (rSig - lowBound) / (highBound - lowBound)
			newG := (gSig - lowBound) / (highBound - lowBound)
			newB := (bSig - lowBound) / (highBound - lowBound)

			img.Pix[i+0] = clampFastFloat(newR * 255.0)
			img.Pix[i+1] = clampFastFloat(newG * 255.0)
			img.Pix[i+2] = clampFastFloat(newB * 255.0)

			if x == 0 && y == 0 {
				fmt.Println(normR, " new ", rSig)
			}
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
