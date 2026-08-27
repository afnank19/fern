package gui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// App is the root application struct, wiring together the image view and sidebar.
type App struct {
	fyneApp   fyne.App
	window    fyne.Window
	imageView *ImageView
	sidebar   *Sidebar
	split     *container.Split
	welcome   *fyne.Container
}

func NewApp() *App {
	a := app.NewWithID("fern-image-processor")
	w := a.NewWindow("Fern")
	a.Settings().SetTheme(newFernTheme())

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

	saveBtn := widget.NewButtonWithIcon(
		"Export Image",
		theme.DocumentSaveIcon(),
		func() {
			saveImageDialog(w, imageView.PrepareExport())
		},
	)

	title := widget.NewLabel("Fern")
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	split := container.NewHSplit(
		imageView.CanvasObject(),
		sidebar.CanvasObject(),
	)
	split.SetOffset(0.75)

	// Welcome overlay: centered icon + text, shown when no image is loaded.
	welcomeIcon := widget.NewIcon(theme.FolderOpenIcon())
	welcomeText := widget.NewLabelWithStyle("Open an image to get started", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	welcome := container.NewCenter(
		container.NewVBox(
			container.NewCenter(welcomeIcon),
			container.NewCenter(welcomeText),
		),
	)
	welcomeBg := container.NewStack(welcome)

	openBtn := widget.NewButtonWithIcon(
		"Open Image",
		theme.FolderOpenIcon(),
		func() {
			openImageDialog(w, func(img *image.RGBA) {
				imageView.LoadImage(img)
				sidebar.Reset()
				welcomeBg.Hide()
				split.Show()
			})
		},
	)

	actions := container.NewHBox(
		openBtn,
		saveBtn,
	)

	statusLabel := imageView.StatusLabel()

	topBar := container.NewBorder(
		nil,          // top
		nil,          // bottom
		title,        // left
		actions,      // right
		statusLabel,  // center
	)

	content := container.NewBorder(
		topBar,
		nil,
		nil,
		nil,
		container.NewMax(welcomeBg, split),
	)

	split.Hide()

	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 750))

	return &App{
		fyneApp:   a,
		window:    w,
		imageView: imageView,
		sidebar:   sidebar,
		split:     split,
		welcome:   welcomeBg,
	}
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
