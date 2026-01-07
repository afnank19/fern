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
	point.PhotoshopGrayscale(rgbaImg)

	image.SaveImage(rgbaImg, "output.png", "./assets/saves")
}
