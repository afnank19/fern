package point

import "image"

func Threshold(img *image.RGBA, threshold uint8) {
	// convert an RGBA image to a GrayScale image for thresholding to work
	bounds := img.Bounds()

	var NewY uint8
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			y := LuminancePhotoshop(r, g, b)
			// using just R, because G,B are the same as well (Grayscale image)
			if y >= threshold {
				NewY = 255
			} else {
				NewY = 0
			}

			img.Pix[i+0] = NewY
			img.Pix[i+1] = NewY
			img.Pix[i+2] = NewY
		}
	}
}
