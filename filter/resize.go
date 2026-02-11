package filter

import (
	"image"
	"math"

	"github.com/afnank19/fern/utils"
)

const downsampleScale = 2.0
const upsampleScale = 2.0

type NgbrPix struct {
	x float64
	y float64
}

func Downsample2x(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	halfWidth := bounds.Dx() / 2
	halfHeight := bounds.Dy() / 2

	out := image.NewRGBA(image.Rect(0, 0, halfWidth, halfHeight))

	outBounds := out.Bounds()
	for y := outBounds.Min.Y; y < outBounds.Max.Y; y++ {
		for x := outBounds.Min.X; x < outBounds.Max.X; x++ {
			// access pixel at (x, y)
			i := out.PixOffset(x, y)

			inputX := (float64(x)+0.5)*downsampleScale - 0.5
			inputY := (float64(y)+0.5)*downsampleScale - 0.5

			dx := inputX - math.Floor(inputX) // horizontal fraction (0 to 1)
			dy := inputY - math.Floor(inputY) // vertical fraction (0 to 1)

			topLeft, topRight, bottomLeft, bottomRight := generateNeighbourhood(inputX, inputY)

			tlIdx := img.PixOffset(int(topLeft.x), int(topLeft.y))
			trIdx := img.PixOffset(int(topRight.x), int(topRight.y))
			blIdx := img.PixOffset(int(bottomLeft.x), int(bottomLeft.y))
			brIdx := img.PixOffset(int(bottomRight.x), int(bottomRight.y))

			out.Pix[i+0] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 0)) // Red
			out.Pix[i+1] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 1)) // Green
			out.Pix[i+2] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 2)) // Blue
			out.Pix[i+3] = 255                                                                                     // Full opaqueness
		}
	}

	return out
}

func generateNeighbourhood(inX, inY float64) (topLeft, topRight, bottomLeft, bottomRight NgbrPix) {
	topLeft.x = math.Floor(inX)
	topLeft.y = math.Floor(inY)

	topRight.x = math.Ceil(inX)
	topRight.y = math.Floor(inY)

	bottomLeft.x = math.Floor(inX)
	bottomLeft.y = math.Ceil(inY)

	bottomRight.x = math.Ceil(inX)
	bottomRight.y = math.Ceil(inY)

	return topLeft, topRight, bottomLeft, bottomRight
}

// Pass 0 for R, 1 for G, 2 for B in channel
func bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx int, img *image.RGBA, dx, dy float64, channel int) float64 {
	return float64(img.Pix[tlIdx+channel])*(1-dx)*(1-dy) +
		float64(img.Pix[trIdx+channel])*dx*(1-dy) +
		float64(img.Pix[blIdx+channel])*(1-dx)*dy +
		float64(img.Pix[brIdx+channel])*dx*dy
}

func Upsample2x(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	doubleWidth := bounds.Dx() * 2
	doubleHeight := bounds.Dy() * 2

	if doubleHeight == 1 || doubleWidth == 1 {
		return img
	}

	out := image.NewRGBA(image.Rect(0, 0, doubleWidth, doubleHeight))

	outBounds := out.Bounds()
	for y := outBounds.Min.Y; y < outBounds.Max.Y; y++ {
		for x := outBounds.Min.X; x < outBounds.Max.X; x++ {
			// access pixel at (x, y)
			i := out.PixOffset(x, y)

			inputX := (float64(x)) / upsampleScale
			inputY := (float64(y)) / upsampleScale

			x0, y0, x1, y1 := upsample4NearestPixels(inputX, inputY, bounds.Dx(), bounds.Dy())

			// X and Y distance
			dx := inputX - float64(x0)
			dy := inputY - float64(y0)

			tlIdx := img.PixOffset(int(x0), int(y0))
			trIdx := img.PixOffset(int(x1), int(y0))
			blIdx := img.PixOffset(int(x0), int(y1))
			brIdx := img.PixOffset(int(x1), int(y1))

			out.Pix[i+0] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 0)) // Red
			out.Pix[i+1] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 1)) // Green
			out.Pix[i+2] = utils.ClampFastFloat(bilinearInterpolation(tlIdx, trIdx, blIdx, brIdx, img, dx, dy, 2)) // Blue
			out.Pix[i+3] = 255                                                                                     // Full opaqueness
		}
	}

	return out
}

// Gets the 4 nearest pixels for the input image needed for the output image
func upsample4NearestPixels(inX, inY float64, w, h int) (int, int, int, int) {
	x0 := int(math.Floor(inX))
	y0 := int(math.Floor(inY))

	x1 := x0 + 1
	y1 := y0 + 1

	// Clamping so we don't overflow the [RGBA RGBA] array
	if x1 >= w {
		x1 = w - 1
	}
	if y1 >= h {
		y1 = h - 1
	}

	return x0, y0, x1, y1
}
