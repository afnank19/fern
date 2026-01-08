package main

import (
	"fmt"

	"github.com/afnank19/fern/image"
	"github.com/afnank19/fern/noise"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	_, rgbaImg := image.LoadImage("./assets/lara.jpeg")

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
	noise.SaltAndPepper(rgbaImg, 0.005)

	image.SaveImage(rgbaImg, "seasoning.png", "./assets/saves")
}
