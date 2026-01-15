// Unsharp Mask for sharpening the image
package filter

import (
	"image"
)

// TODO: possibly handle blurAmt being too large
func UnsharpMask(img *image.RGBA, blurAmt float64, sharpeningAmount float64 ) {
	newImg := GaussianBlur(img, blurAmt)

	bounds := newImg.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := newImg.PixOffset(x, y)
			r := newImg.Pix[i+0]
			g := newImg.Pix[i+1]
			b := newImg.Pix[i+2]

			rOld := img.Pix[i+0]
			gOld := img.Pix[i+1]
			bOld := img.Pix[i+2]

			rDetail := calculateDetail(rOld, r)
			gDetail := calculateDetail(gOld, g)
			bDetail := calculateDetail(bOld, b)

			img.Pix[i+0] = calculateUnsharpMask(rOld, sharpeningAmount, rDetail)
			img.Pix[i+1] = calculateUnsharpMask(gOld, sharpeningAmount, gDetail)
			img.Pix[i+2] = calculateUnsharpMask(bOld, sharpeningAmount, bDetail)
		}
	}
}


func calculateDetail(original uint8, blurred uint8) int {
    // Return signed integer to preserve negative values
    return int(original) - int(blurred)
}

func calculateUnsharpMask(original uint8, amount float64, detail int) uint8{
	return clampFastFloat(float64(original) + amount * float64(detail))
}
