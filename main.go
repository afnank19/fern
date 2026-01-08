package main

import (
	"fmt"

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
	// point.PhotoshopGrayscale(rgbaImg)
	// point.LinearContrast(rgbaImg, 1.5)
	// point.FastGrayscale(rgbaImg)
	// point.SigmoidalContrast(rgbaImg, 0.1)
	point.Threshold(rgbaImg, 70)

	image.SaveImage(rgbaImg, "output.png", "./assets/saves")
}
