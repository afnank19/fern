package filter

import (
	"image"
)

// sharpening does not depend on kernel size, but rather strength
const KERNEL_SIZE = 3

func Sharpen(img *image.RGBA, val float64) *image.RGBA{
	kernel := generateSharpenKernel(val)

	return Convolve(img, kernel)
}

// Works for a 3x3 only
func generateSharpenKernel(val float64) [][]float64{
	kernel := make([][]float64, KERNEL_SIZE)
	for i := range kernel {
		kernel[i] = make([]float64, KERNEL_SIZE)
	}

	row, col, _ := findMiddleOfKernel(kernel)

	for i := range kernel {
		for j := range kernel {
			if i == 0 && (j == 0 || j == KERNEL_SIZE-1) {
				kernel[i][j] = 0
				continue
			}

			if i == KERNEL_SIZE-1 && (j == 0 || j == KERNEL_SIZE-1) {
				kernel[i][j] = 0
				continue
			}

			if i == row && col == j {
				kernel[i][j] = 1 + 4 * val
				continue
			}

			kernel[i][j] = -val
		}
	}

	return kernel;
}
