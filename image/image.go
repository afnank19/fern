package image

import (
	"fmt"
	"image"
	"image/color"
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

func IterateImage(img *image.RGBA) {
	bounds := img.Bounds()

	fmt.Println()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			if x > 2 {
				break
			}
			fmt.Print("[", x, y, "] ")

		}
		fmt.Println()
		if y > 2 {
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
		panic(err)
	}
	defer out.Close()

	png.Encode(out, img)
}

func CopyRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func TestNegAccess(img *image.RGBA) {
	i := img.PixOffset(-1, -1)
	fmt.Println(i)
}

// TODO: Remove or reuse properly
// Was to generate a small photo so i could analyze pixels
func SaveBlah() {
	out := image.NewRGBA(image.Rect(0, 0, 3, 3))

	// Loop through every pixel
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			// Calculate a red value that changes based on position
			// Using (x + y * 4) * 16 gives us 16 distinct shades of red
			r := uint8((x + y*4) * 16)

			// Set the pixel color: RGBA(red, green, blue, alpha)
			// Alpha is set to 255 for full opacity
			out.SetRGBA(x, y, color.RGBA{R: r, G: 0, B: 0, A: 255})
		}
	}

	SaveImage(out, "sample4x4.png", "./assets/samples")
}
