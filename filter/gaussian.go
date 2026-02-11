package filter

import (
	"fmt"
	"image"
	"math"
)

const MAX_BLUR_STRENGTH float64 =  0.03

func GaussianBlur(img *image.RGBA, val float64) *image.RGBA {

	kernelSize, sigma := calculateGaussianKernelSize(img, val)
	kernel := buildGaussianKernel(kernelSize, sigma)

	return Convolve(img, kernel)
}

func buildGaussianKernel(kernelSize int, sigma float64) [][]float64{
	fmt.Println("building gaussian kernel with size", kernelSize, " and sigma", sigma)
	kernel := make([][]float64, kernelSize)
	for i := range kernel {
		kernel[i] = make([]float64, kernelSize)
	}

	radius := kernelSize / 2

	var kernelValSum float64= 0.0
	for i := range kernel {
		for j := range kernel[i] {
			x := float64(i - radius)
			y := float64(j - radius)

			value := math.Exp( -(x*x + y*y)/ (2 * (sigma * sigma)))
			kernel[i][j] = value
			kernelValSum += value
		}
	}

	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] = kernel[i][j] /  kernelValSum
			// fmt.Print(kernel[i][j], " ")
		}
		// fmt.Println()
	}

	return kernel
}

//NOTE: Add a safety cap (optional)?
func calculateGaussianKernelSize(img *image.RGBA, val float64) (int ,float64){
	b := img.Bounds()

	sigmaMax := MAX_BLUR_STRENGTH * math.Min(float64(b.Dx()), float64(b.Dy()))

	// this gives finer control over low blur values, whereas just have val does not, use if needed
	sigma := (val * val) * sigmaMax

	// sigma := val * sigmaMax
 	if sigma == 0 {
        return 1, 0.0
    }

 	radius := int(math.Ceil(3.0 * sigma))

    return 2*radius + 1, sigma
}


func FastGaussianBlur(img *image.RGBA, val float64)  *image.RGBA{
	kernelSize, sigma := calculateGaussianKernelSize(img, val)
	kernel := build1DGaussianKernel(kernelSize, sigma)

	return Convolve1D(img, kernel)
}

func build1DGaussianKernel(kernelSize int, sigma float64) []float64{
	fmt.Println("building 1d gaussian kernel with size", kernelSize, " and sigma", sigma)
	kernel := make([]float64, kernelSize)

	radius := kernelSize / 2

	var kernelValSum float64= 0.0
	for i := range kernel {
			x := float64(i - radius)

			value := math.Exp( -(x*x)/ (2 * (sigma * sigma)))
			kernel[i] = value
			kernelValSum += value
	}

	for i := range kernel {
		kernel[i] = kernel[i] /  kernelValSum
	}

	return kernel
}

func GetKernelForBloom(img *image.RGBA, val float64) []float64 {
	kernelSize, sigma := calculateGaussianKernelSize(img, val)
	kernel := build1DGaussianKernel(kernelSize, sigma)

	return kernel
}

func FastGaussianBlurWithKernel(img *image.RGBA, kernel []float64) *image.RGBA {
	return Convolve1D(img, kernel)
}
