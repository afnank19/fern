package main

import (
	"fmt"

	"github.com/afnank19/fern/internal/gui"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	// _, rgbaImg := image.LoadImage("./assets/samples/harris.JPG")

	// composite.Bloom(rgbaImg, 0.75, 0.9, 0.85)
	// geometric.ChromaticAberration(rgbaImg, 2, 0)
	// noise.Gaussian(rgbaImg, 30, true)

	// image.SaveImage(rgbaImg, "bloom.png", "./assets/bloom")
	app := gui.NewApp()
	app.Run()
}
