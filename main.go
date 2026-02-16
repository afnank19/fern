package main

import (
	"fmt"

	"github.com/afnank19/fern/internal/gui"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	// _, rgbaImg := image.LoadImage("./assets/samples/trail.JPG")

	// composite.Bloom(rgbaImg, 0.21, 0.48, 0.35)
	// // noise.Gaussian(rgbaImg, 10, true)
	// geometric.ChromaticAberration(rgbaImg, 2, 0)

	// image.SaveImage(rgbaImg, "bloom.png", "./assets/bloom")
	app := gui.NewApp()
	app.Run()
}
