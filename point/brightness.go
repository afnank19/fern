package point

import "image"

func Brightness(img *image.RGBA, value int) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			newR := int(r) + (value)
			newG := int(g) + (value)
			newB := int(b) + (value)

			img.Pix[i+0] = clamp(newR)
			img.Pix[i+1] = clamp(newG)
			img.Pix[i+2] = clamp(newB)
		}
	}
}

// may move to a utils later
func clamp(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
