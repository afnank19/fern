package noise

import (
	"fmt"
	"image"
	"math/rand"
)

func SaltAndPepper(img *image.RGBA, p float64) {
	bounds := img.Bounds()

	salt_p := p / 2.0
	fmt.Println("Salt Probability", salt_p)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			sample := rand.Float64()

			if sample < salt_p {
				img.Pix[i+0] = 255
				img.Pix[i+1] = 255
				img.Pix[i+2] = 255
			} else if sample < p {
				img.Pix[i+0] = 0
				img.Pix[i+1] = 0
				img.Pix[i+2] = 0
			}
		}
	}
}
