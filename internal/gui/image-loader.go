package gui

import (
	"image"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// openImageDialog presents a file picker and calls onLoaded with the decoded
// image on success. Decoding is kept here so the rest of the GUI stays clean.
func openImageDialog(w fyne.Window, onLoaded func(*image.RGBA)) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		decoded, _, err := image.Decode(reader)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Normalise to *image.RGBA so processing functions always get the same type.
		// TODO: Compare with LoadImage from local img pkg
		bounds := decoded.Bounds()
		rgba := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rgba.Set(x, y, decoded.At(x, y))
			}
		}

		onLoaded(rgba)
	}, w)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}

// decodeImage decodes using standard image.Decode. reader is an io.Reader.
// We separate this to keep error handling / tests easy.
func decodeImage(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	return img, format, err
}
