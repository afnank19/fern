package gui

import (
	"image"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// openImageDialog creates a file-open dialog, decodes the image, and updates the canvas.Image.
func openImageDialog(parent fyne.Window, imgCanvas *canvas.Image) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if reader == nil { // user cancelled
			return
		}
		defer reader.Close()

		decoded, _, decodeErr := decodeImage(reader)
		if decodeErr != nil {
			dialog.ShowError(decodeErr, parent)
			return
		}

		// Update canvas image and refresh UI.
		// Assign the decoded image (image.Image) and refresh to show it:
		imgCanvas.Image = decoded
		imgCanvas.FillMode = canvas.ImageFillContain
		imgCanvas.Refresh()
	}, parent)

	// Limit selectable files to common image extensions
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif", ".webp"}))
	fd.Show()
}

// decodeImage decodes using standard image.Decode. reader is an io.Reader.
// We separate this to keep error handling / tests easy.
func decodeImage(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	return img, format, err
}
