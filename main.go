package main

import (
	"fmt"
	"time"

	"github.com/afnank19/fern/image"
	"github.com/afnank19/fern/point"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	_, rgbaImg := image.LoadImage("./assets/evelyn.png")

	// image.IterateImage(testImage)
	// point.Invert(rgbaImg)
	// point.Brightness(rgbaImg, 40)
	// point.Grayscale(rgbaImg)
	// point.AvgGrayscale(rgbaImg)
	///point.PhotoshopGrayscale(rgbaImg)
	// point.LinearContrast(rgbaImg, 1.5)
	// point.FastGrayscale(rgbaImg)
	// point.SigmoidalContrast(rgbaImg, 0.10)
	start := time.Now()
	point.FastSigmoidalContrast(rgbaImg, 8.0)
	elapsed := time.Since(start)
	fmt.Println(" Elapsed -", elapsed)
	// point.Threshold(rgbaImg, 50)

	// noise.Uniform(rgbaImg, 20, true)
	// noise.Gaussian(rgbaImg, 20, true)
	// noise.SaltAndPepper(rgbaImg, 0.005)

	// filter.SliceShenanas()
	// image.IterateImage(rgbaImg)
	// image.TestNegAccess(rgbaImg)
	// newImg := filter.BoxBlur(rgbaImg, 1.0)
	// fmt.Println(filter.CalculateKernelSize(rgbaImg, 1.0))

	image.SaveImage(rgbaImg, "normalsigcontra.png", "./assets/saves")

	// image.SaveImage(newImg, "filter.png", "./assets/saves")
}
