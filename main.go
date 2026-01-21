package main

import (
	"fmt"

	"github.com/afnank19/fern/composite"
	"github.com/afnank19/fern/image"
	"github.com/afnank19/fern/noise"
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
	// start := time.Now()
	// point.FastSigmoidalContrast(rgbaImg, 8.0)
	// elapsed := time.Since(start)
	// fmt.Println(" Elapsed -", elapsed)
	// point.Threshold(rgbaImg, 50)

	// noise.Uniform(rgbaImg, 20, true)
	// noise.Gaussian(rgbaImg, 20, true)
	// noise.SaltAndPepper(rgbaImg, 0.005)

	// filter.SliceShenanas()
	// image.IterateImage(rgbaImg)
	// image.TestNegAccess(rgbaImg)
	// newImg := filter.BoxBlur(rgbaImg, 1.0)
	// fmt.Println(filter.CalculateKernelSize(rgbaImg, 1.0))
	//
	// start := time.Now()
	// newImg := filter.GaussianBlur(rgbaImg, 0.01)
	// elapsed := time.Since(start)
	// fmt.Println(" Elapsed -", elapsed)

	// newImg := filter.Sharpen(rgbaImg, 0.5)

	// filter.UnsharpMask(rgbaImg, 0.2, 1.5)
	composite.NaiveBloom(rgbaImg, 0.75, 0.67, 0.9)
	noise.Gaussian(rgbaImg, 5, true)

	image.SaveImage(rgbaImg, "bloom.png", "./assets/saves")

	// image.SaveImage(newImg, "sharpen.png", "./assets/saves")
}
