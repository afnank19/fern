package noise

import (
	"image"
	"math/rand"
)

func Uniform(img *image.RGBA, k int, perChannel bool) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			noiseR := rand.Intn(2*k+1) - k
			noiseG := noiseR
			noiseB := noiseR

			if perChannel {
				noiseG = rand.Intn(2*k+1) - k
				noiseB = rand.Intn(2*k+1) - k
			}

			// noise RGB will have same values unless perChannel is enabled
			newR := clamp(int(r) + noiseR)
			newG := clamp(int(g) + noiseG)
			newB := clamp(int(b) + noiseB)

			img.Pix[i+0] = newR
			img.Pix[i+1] = newG
			img.Pix[i+2] = newB
		}
	}
}

func clamp(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
