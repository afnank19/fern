package filter

import (
	"image"
	"math"
)

// TODO: Move into box blur file
// TODO: Gaussian Blur
func BoxBlur(img *image.RGBA, val float64) *image.RGBA {
	if val > 1.0 {
		return img
	}

	kernelSize := CalculateKernelSize(img, val)
	kernel := buildAvgKernel(kernelSize)

	return Convolve(img, kernel)
}

func buildAvgKernel(kernelSize int) [][]float64 {
	kernel := make([][]float64, kernelSize)
	for i := range kernel {
		kernel[i] = make([]float64, kernelSize)
	}

	kernelValue := 1.0 / float64(kernelSize*kernelSize)

	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] = kernelValue
		}
	}

	return kernel
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
