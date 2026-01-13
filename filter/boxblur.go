package filter

import (
	"image"
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
