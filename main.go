package main

import (
	"fmt"

	"github.com/afnank19/fern/image"
	"github.com/afnank19/fern/noise"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	_, rgbaImg := image.LoadImage("./assets/saves/bloom.png")
	// _, rgbaImg := image.LoadImage("./assets/samples/se98.jpg")

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
	// composite.NaiveBloom(rgbaImg, 0.30, 0.30, 0.8)
	noise.Gaussian(rgbaImg, 3, true)
	// noise.Uniform(rgbaImg, 5, true)

	// newImg := filter.Downsample2x(rgbaImg)
	// newImg = filter.Downsample2x(newImg)

	image.SaveImage(rgbaImg, "fern.png", "./assets/saves")

	// start := time.Now()
	// newImg := filter.GaussianBlur(rgbaImg, 1.5)
	// elapsed := time.Since(start)
	// fmt.Println(" NAIVE GAUSSIAN: Elapsed -", elapsed)

	// image.SaveImage(newImg, "naive-gauss.png", "./assets/saves")

	fmt.Println("----------------------------")

	// start = time.Now()
	// fastGauss := filter.FastGaussianBlur(rgbaImg, 1.5)
	// elapsed = time.Since(start)
	// fmt.Println(" FAST GAUSSIAN: Elapsed -", elapsed)

	// image.SaveImage(fastGauss, "fast-gauss.png", "./assets/saves")
}
