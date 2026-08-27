package image

import (
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
