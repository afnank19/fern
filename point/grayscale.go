package point

import (
	"image"
	"math"
)

// Slower implementation, uses float arithmetic
func Grayscale(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			floatY := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)

			y := uint8(math.Round(floatY))

			img.Pix[i+0] = y
			img.Pix[i+1] = y
			img.Pix[i+2] = y
		}
	}
}

func FastGrayscale(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			yInt32 := (13933*uint32(r) + 46871*uint32(g) + 4732*uint32(b)) >> 16

			y := uint8(yInt32)

			img.Pix[i+0] = y
			img.Pix[i+1] = y
			img.Pix[i+2] = y
		}
	}
}

// Uses the most simplest and basic approach
func AvgGrayscale(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			floatY := (float64(r) + float64(g) + float64(b)) / 3.0

			y := uint8(math.Round(floatY))

			img.Pix[i+0] = y
			img.Pix[i+1] = y
			img.Pix[i+2] = y
		}
	}
}

// Corrects the gamma by converting to pixel vals to linear before computing the grayscale
func PhotoshopGrayscale(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			// rLin := sRGBToLinear(r)
			// gLin := sRGBToLinear(g)
			// bLin := sRGBToLinear(b)

			// yLin := (0.2126 * rLin) + (0.7152 * gLin) + (0.0722 * bLin) // weighted formula for a pixel avg

			// y := linearToSRGB(yLin)
			y := LuminancePhotoshop(r, g, b)

			img.Pix[i+0] = y
			img.Pix[i+1] = y
			img.Pix[i+2] = y
		}
	}
}

func LuminancePhotoshop(r, g, b uint8) uint8 {
	rLin := sRGBToLinear(r)
	gLin := sRGBToLinear(g)
	bLin := sRGBToLinear(b)

	yLin := (0.2126 * rLin) + (0.7152 * gLin) + (0.0722 * bLin) // weighted formula for a pixel avg

	y := linearToSRGB(yLin)

	return y
}

// Convert a standard sRGB uint8 value to a Linear float64 (0.0 to 1.0)
func sRGBToLinear(c uint8) float64 {
	val := float64(c) / 255.0
	if val <= 0.04045 {
		return val / 12.92
	}
	return math.Pow((val+0.055)/1.055, 2.4)
}

// Convert a Linear float64 back to a standard sRGB uint8
func linearToSRGB(val float64) uint8 {
	var srgb float64
	if val <= 0.0031308 {
		srgb = val * 12.92
	} else {
		srgb = 1.055*math.Pow(val, 1.0/2.4) - 0.055
	}
	return uint8(math.Round(srgb * 255.0))
}
