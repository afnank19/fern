package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	App    fyne.App
	Window fyne.Window

	imageCanvas *canvas.Image
}

func BuildGui() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}

func NewApp() *App {
	a := app.NewWithID("fern-image-processor")
	w := a.NewWindow("Hello World")
	// empty canvas image to start with
	imgCanvas := canvas.NewImageFromImage(nil)
	imgCanvas.FillMode = canvas.ImageFillContain // keep aspect ratio

	// top bar with an Open button
	openBtn := widget.NewButton("Open...", func() {
		openImageDialog(w, imgCanvas)
	})

	topBar := container.NewHBox(openBtn)

	content := container.NewBorder(topBar, nil, nil, nil, container.NewMax(imgCanvas))
	w.SetContent(content)
	w.Resize(fyne.NewSize(1000, 700))

	return &App{
		App:         a,
		Window:      w,
		imageCanvas: imgCanvas,
	}
}

func (app *App) Run() {
	app.Window.ShowAndRun()
}
