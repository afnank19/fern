package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type App struct {
    app fyne.App
    Window  fyne.Window
}

func BuildGui() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}

func NewApp() *App {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetContent(widget.NewLabel("Hello World!"))

	return &App{
		app: a,
		Window: w,
	}
}

func (app *App) Run() {
	app.Window.ShowAndRun()
}
