package main

import (
	"fmt"

	"github.com/afnank19/fern/filter"
	"github.com/afnank19/fern/image"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	_, rgbaImg := image.LoadImage("./assets/what.jpeg")

	// image.IterateImage(testImage)
	// point.Invert(rgbaImg)
	// point.Brightness(rgbaImg, 40)
	// point.Grayscale(rgbaImg)
	// point.AvgGrayscale(rgbaImg)
	// point.PhotoshopGrayscale(rgbaImg)
	// point.LinearContrast(rgbaImg, 1.5)
	// point.FastGrayscale(rgbaImg)
	// point.SigmoidalContrast(rgbaImg, 0.1)
	// point.Threshold(rgbaImg, 150)

	// noise.Uniform(rgbaImg, 20, true)
	// noise.Gaussian(rgbaImg, 20, true)
	// noise.SaltAndPepper(rgbaImg, 0.005)

	// filter.SliceShenanas()
	// image.IterateImage(rgbaImg)
	// image.TestNegAccess(rgbaImg)
	newImg := filter.BoxBlur(rgbaImg, 0.1)
	// fmt.Println(filter.CalculateKernelSize(rgbaImg, 1.0))

	// image.SaveImage(rgbaImg, "seasoning.png", "./assets/saves")

	image.SaveImage(newImg, "filter.png", "./assets/saves")
}
