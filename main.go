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
	point.Brightness(rgbaImg, 40)

	image.SaveImage(rgbaImg, "brightness.png", "./assets/saves")
}
