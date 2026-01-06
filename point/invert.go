package point

import (
	"image"
)

func Invert(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			img.Pix[i+0] = 255 - r
			img.Pix[i+1] = 255 - g
			img.Pix[i+2] = 255 - b
		}
	}
}
