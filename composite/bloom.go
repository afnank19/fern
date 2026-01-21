package composite

import (
	"image"
	"math"

	"github.com/afnank19/fern/filter"
	localImg "github.com/afnank19/fern/image"
	"github.com/afnank19/fern/point"
	"github.com/afnank19/fern/utils"
)

func NaiveBloom(img *image.RGBA, intensity ,threshold, blurAmt float64) {
	bounds := img.Bounds()

	brightPass := image.NewRGBA(bounds)

	// Create an image that has the bright parts based on the threshold
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]
			a := img.Pix[i+3]

			lum := point.LuminancePhotoshop(r, g, b)
			// IMPORTANT: convert either on to the other properly, not just type conv
			if (lum < utils.FloatToUint8(threshold)) {
				// fmt.Println("lum", lum, "thre", utils.FloatToUint8(threshold))
				// Softer bloom threshold
				t := float64(lum)/255.0
				w := math.Max(0, (t-threshold)/(1-threshold))

				r = uint8(float64(r) * w)
				g = uint8(float64(g) * w)
				b = uint8(float64(b) * w)
			}

			// brightIndex := brightPass.PixOffset(x, y)
			brightPass.Pix[i + 0] = r
			brightPass.Pix[i + 1] = g
			brightPass.Pix[i + 2] = b
			brightPass.Pix[i + 3] = a
		}
	}


	localImg.SaveImage(brightPass, "brightpass.png", "./assets/saves")
	blurredBright := filter.GaussianBlur(brightPass, blurAmt)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			img.Pix[i+0] = addBloomToPixel(r, intensity, blurredBright.Pix[i+0])
			img.Pix[i+1] = addBloomToPixel(g, intensity, blurredBright.Pix[i+1])
			img.Pix[i+2] = addBloomToPixel(b, intensity, blurredBright.Pix[i+2])
		}
	}

}

func addBloomToPixel(original uint8, intensity float64, blurredHighlights uint8) uint8 {
	return utils.ClampFastFloat(float64(original) + float64(blurredHighlights) * intensity)
}
