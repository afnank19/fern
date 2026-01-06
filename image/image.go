package image

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"os"
)

// may stick with image.RGBA and not image.Image
func LoadImage(path string) (image.Image, *image.RGBA) {
	imgFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)

	if err != nil {
		panic(err)
	}

	// convert any format into RGBA
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src) // mutates the img var

	return img, rgba
}

func IterateImage(img image.Image) {
	bounds := img.Bounds()

	fmt.Println()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			fmt.Print("X:", x, y, img.At(x, y), " ")
			if x > 0 {
				break
			}

		}
		fmt.Println()
		if y > 0 {
			break
		}
	}

	// r, g, b, a := img.At(0, 0).RGBA()
	// fmt.Println("RGBA", r>>8, g>>8, b>>8, a>>8)
}

// If path is an empty string, save in /assets/saves
func SaveImage(img *image.RGBA, name, path string) {
	out, err := os.Create(path + "/" + name)
	if err != nil {
		// handle error
	}
	defer out.Close()

	png.Encode(out, img)
}
