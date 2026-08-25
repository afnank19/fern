package gui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// App is the root application struct, wiring together the image view and sidebar.
type App struct {
	fyneApp   fyne.App
	window    fyne.Window
	imageView *ImageView
	sidebar   *Sidebar
}

func NewApp() *App {
	a := app.NewWithID("fern-image-processor")
	w := a.NewWindow("Fern")

	imageView := NewImageView()

	var sidebar *Sidebar
	sidebar = NewSidebar(SidebarCallbacks{
		OnParam: imageView.SetParam,
		OnCommit: func(kind OpKind) {
			imageView.Commit(kind)
			sidebar.ResetOp(kind)
		},
		OnUndo: imageView.Undo,
	})

	openBtn := widget.NewButton("Open Image", func() {
		openImageDialog(w, func(img *image.RGBA) {
			imageView.LoadImage(img)
			sidebar.Reset()
		})
	})

	saveBtn := widget.NewButton("Export Image", func() {
		saveImageDialog(w, imageView.PrepareExport())
	})

	title := widget.NewLabel("Retarted Photoshop")

	topBar := container.NewHBox(openBtn, saveBtn, title)

	split := container.NewHSplit(
		imageView.CanvasObject(),
		sidebar.CanvasObject(),
	)
	split.SetOffset(0.75) // 75% image, 25% sidebar

	content := container.NewBorder(topBar, nil, nil, nil, split)
	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 750))

	return &App{
		fyneApp:   a,
		window:    w,
		imageView: imageView,
		sidebar:   sidebar,
	}
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
