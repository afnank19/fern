package point

import "image"

func Threshold(img *image.RGBA, threshold uint8) {
	// convert an RGBA image to a GrayScale image for thresholding to work
	PhotoshopGrayscale(img)
	// Possible suggestion: Process the pixel to be grayscale, and then calculate the threshold on it
	// This makes it be one loop, instead of going over the pixels once for grayscale, and then again for threshold
	// but first remove repeated code in other implementations so that they can be used here!

	bounds := img.Bounds()

	var NewY uint8
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]

			// using just R, because G,B are the same as well (Grayscale image)
			if r >= threshold {
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
