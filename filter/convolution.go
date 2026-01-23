package filter

import (
	"errors"
	"image"
	"math"

	localImg "github.com/afnank19/fern/image"
)

var AVG_KERNEL = [][]float64{
	{0.11111111, 0.11111111, 0.11111111},
	{0.11111111, 0.11111111, 0.11111111},
	{0.11111111, 0.11111111, 0.11111111},
}

// 7x7 average (box blur) kernel
var Avg7x7 = [][]float64{
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
	{0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408, 0.020408},
}

var TEST_KERNEL = [][]float64{
	{0.1, 0.2, 0.3},
	{0.4, 0.5, 0.6},
	{0.7, 0.8, 0.9},
}

// always work on the "working" image, and never the source img
func Convolve(img *image.RGBA, kernel [][]float64) *image.RGBA {
	workingImg := localImg.CopyRGBA(img)

	bounds := img.Bounds()

	midX, midY, _ := findMiddleOfKernel(kernel)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			xD := x - midX
			yD := y - midY

			sumRGB := applyKernelOnPixel(xD, yD, bounds, img, kernel)

			i := workingImg.PixOffset(x, y)
			workingImg.Pix[i+0] = clampFastFloat(sumRGB[0])
			workingImg.Pix[i+1] = clampFastFloat(sumRGB[1])
			workingImg.Pix[i+2] = clampFastFloat(sumRGB[2])
		}
	}

	return workingImg
}

func applyKernelOnPixel(x, y int, bounds image.Rectangle, img *image.RGBA, kernel [][]float64) []float64 {
	xCopy := x

	var sumR float64 = 0.0
	var sumG float64 = 0.0
	var sumB float64 = 0.0

	for i := 0; i < len(kernel); i++ {
		for j := 0; j < len(kernel); j++ {
			// fmt.Print("[", x, y, "]", " ")
			kernelVal := kernel[i][j]
			// fmt.Print("[", kernelVal, "- ", x, y, "]", " ")
			//
			if x < bounds.Min.X || y < bounds.Min.Y || x >= bounds.Max.X || y >= bounds.Max.Y {
				x++
				continue
			} else {
				i := img.PixOffset(x, y)
				r := img.Pix[i+0]
				g := img.Pix[i+1]
				b := img.Pix[i+2]

				tempR := kernelVal * float64(r)
				tempG := kernelVal * float64(g)
				tempB := kernelVal * float64(b)

				sumR += tempR
				sumG += tempG
				sumB += tempB
				x++
			}
		}
		// fmt.Println()
		y++
		x = xCopy
	}

	return []float64{sumR, sumG, sumB}
}

func findMiddleOfKernel(data [][]float64) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, errors.New("empty slice")
	}
	if len(data[0]) == 0 {
		return 0, 0, errors.New("empty row")
	}

	midRow := len(data) / 2
	midCol := len(data[0]) / 2

	return midRow, midCol, nil
}

func clampFastFloat(val float64) uint8 {
	// This achieves the same result as the if-statements
	return uint8(math.Max(0, math.Min(255, val)))
}

// Calculates the kernel size based on image width
// returns the kernel width, e.g a width of 3 means a 3x3 kernel
func CalculateKernelSize(img *image.RGBA, val float64) int {
	b := img.Bounds()

	mappedVal := val * 0.025

	kernelSizeF := (2. * (mappedVal * math.Min(float64(b.Max.X), float64(b.Max.Y)))) + 1.

	kernelSize := int(math.Ceil(kernelSizeF))

	if kernelSize%2 == 0 {
		kernelSize-- // could ++, but eh
	}

	return kernelSize
}

func Convolve1D(img *image.RGBA, kernel []float64) *image.RGBA {
	workingImg := localImg.CopyRGBA(img)

	bounds := img.Bounds()

	midX, _ := findMiddleOf1DKernel(kernel)

	// Horizontal Pass
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			xD := x - midX

			sumRGB := apply1DKernelOnPixel(xD, y, bounds, img, kernel, false)

			i := workingImg.PixOffset(x, y)
			workingImg.Pix[i+0] = clampFastFloat(sumRGB[0])
			workingImg.Pix[i+1] = clampFastFloat(sumRGB[1])
			workingImg.Pix[i+2] = clampFastFloat(sumRGB[2])
		}
	}

	out := localImg.CopyRGBA(workingImg)
	// Vertical Pass on horizontal pass result
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			yD := y - midX

			sumRGB := apply1DKernelOnPixel(x, yD, bounds, workingImg, kernel, true)

			i := out.PixOffset(x, y)
			out.Pix[i+0] = clampFastFloat(sumRGB[0])
			out.Pix[i+1] = clampFastFloat(sumRGB[1])
			out.Pix[i+2] = clampFastFloat(sumRGB[2])
		}
	}

	return out
}

func findMiddleOf1DKernel(data []float64) (int, error) {
	if len(data) == 0 {
		return 0, errors.New("empty slice")
	}

	midRow := len(data) / 2

	return midRow, nil
}

func apply1DKernelOnPixel(x, y int, bounds image.Rectangle, img *image.RGBA, kernel []float64, vertical bool) []float64 {
	var sumR float64 = 0.0
	var sumG float64 = 0.0
	var sumB float64 = 0.0

	// Convolve vertically
	if vertical {
		for idx := 0; idx < len(kernel); idx++ {
			// Skip if y is out of bounds (Same as multiplying kernel value with 0)
			if y < bounds.Min.Y || y >= bounds.Max.Y {
				y++
				continue
			} else {
				i := img.PixOffset(x, y)
				r := img.Pix[i+0]
				g := img.Pix[i+1]
				b := img.Pix[i+2]

				sumR += kernel[idx] * float64(r)
				sumG += kernel[idx] * float64(g)
				sumB += kernel[idx] * float64(b)
				y++
			}
		}

		return []float64{sumR, sumG, sumB}
	}


	for idx := 0; idx < len(kernel); idx++ {
		// Skip if x is out of bounds (Same as multiplying kernel value with 0)
		if x < bounds.Min.X || x >= bounds.Max.X {
			x++
			continue
		} else {
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			sumR += kernel[idx] * float64(r)
			sumG += kernel[idx] * float64(g)
			sumB += kernel[idx] * float64(b)
			x++
		}
	}

	return []float64{sumR, sumG, sumB}
}
